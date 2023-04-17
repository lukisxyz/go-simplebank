package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocks "github.com/flukis/simplebank/db/mock"
	db "github.com/flukis/simplebank/db/sqlc"
	"github.com/flukis/simplebank/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetchAccountAPI(t *testing.T) {
	n := 5

	account := make([]db.Account, n)
	for i := 0; i < n; i++ {
		account[i] = randomAccount()
	}

	type falseFetchAccountRequest struct {
		PageID int32  `form:"page" binding:"required"`
		Limit  string `form:"limit" binding:"required"`
	}

	testCases := []struct {
		name  string
		body  any
		build func(store *mocks.Store)
		check func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			body: fetchAccountRequest{
				PageID: 1,
				Limit:  int32(n),
			},
			build: func(store *mocks.Store) {
				arg := db.FetchAccountsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.On("FetchAccounts", mock.Anything, arg).
					Return(account, nil).
					Once()
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchAccount(t, rec.Body, fetchAccountSuccessResponse{
					Data: account,
					Meta: Meta{
						Limit: int32(n),
						Page:  1,
					},
				})
			},
		},
		{
			name: "StatusBadRequestNotValid",
			body: fetchAccountRequest{
				PageID: 0,
				Limit:  int32(n),
			},
			build: func(store *mocks.Store) {
				arg := db.FetchAccountsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.On("FetchAccounts", mock.Anything, arg).
					Return(account, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusBadRequestWrongParamType",
			body: falseFetchAccountRequest{
				PageID: 0,
				Limit:  "w",
			},
			build: func(store *mocks.Store) {
				arg := db.FetchAccountsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.On("FetchAccounts", mock.Anything, arg).
					Return(account, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusInternalServerError",
			body: fetchAccountRequest{
				PageID: 1,
				Limit:  int32(n),
			},
			build: func(store *mocks.Store) {
				arg := db.FetchAccountsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.On("FetchAccounts", mock.Anything, arg).
					Return(account, sql.ErrConnDone)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "StatusNotFound1",
			body: fetchAccountRequest{
				PageID: 1,
				Limit:  int32(n),
			},
			build: func(store *mocks.Store) {
				arg := db.FetchAccountsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.On("FetchAccounts", mock.Anything, arg).
					Return(make([]db.Account, 5), sql.ErrNoRows)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			name: "StatusNotFound2",
			body: fetchAccountRequest{
				PageID: 1,
				Limit:  int32(5),
			},
			build: func(store *mocks.Store) {
				arg := db.FetchAccountsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.On("FetchAccounts", mock.Anything, arg).
					Return(make([]db.Account, 0), nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
	}

	for _, ts := range testCases {
		t.Run(ts.name, func(t *testing.T) {
			store := &mocks.Store{}
			ts.build(store)

			dur, err := time.ParseDuration("1m")
			require.NoError(t, err)

			server, err := NewServer(store, util.Config{
				TokenSymetricKey:    "12345678901234567890123456789012",
				AccessTokenDuration: dur,
			})
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			data, err := json.Marshal(ts.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodGet, "/account", bytes.NewReader(data))
			require.NoError(t, err)
			req.Header = http.Header{
				"Content-Type": {"application/json"},
			}

			server.router.ServeHTTP(rec, req)
			ts.check(t, rec)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()

	type falseCreateAccountRequest struct {
		OwnerID  uuid.UUID `json:"owner_id" binding:"required"`
		Currency string    `json:"currency" binding:"required,oneof=USD EUR"`
		Balance  string    `json:"balance" binding:"requied"`
	}

	testCases := []struct {
		name  string
		body  any
		build func(store *mocks.Store)
		check func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			body: createAccountRequest{
				OwnerID:  account.OwnerID,
				Currency: account.Currency,
				Balance:  account.Balance,
			},
			build: func(store *mocks.Store) {
				arg := db.CreateAccountParams{
					OwnerID:  account.OwnerID,
					Currency: account.Currency,
					Balance:  account.Balance,
				}
				store.On("CreateAccount", mock.Anything, arg).
					Return(account, nil).
					Once()
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchAccount(t, rec.Body, createAccountSuccessResponse{Data: account})
			},
		},
		{
			name: "StatusBadRequestNotValidParams",
			body: createAccountRequest{
				OwnerID:  account.OwnerID,
				Currency: "PESO",
				Balance:  account.Balance,
			},
			build: func(store *mocks.Store) {
				arg := db.CreateAccountParams{
					OwnerID:  account.OwnerID,
					Currency: account.Currency,
					Balance:  account.Balance,
				}
				store.On("CreateAccount", mock.Anything, arg).
					Return(db.Account{}, mock.Anything)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusBadRequestWrongParams",
			body: falseCreateAccountRequest{
				OwnerID:  account.OwnerID,
				Currency: account.Currency,
				Balance:  "9899",
			},
			build: func(store *mocks.Store) {
				arg := db.CreateAccountParams{
					OwnerID:  account.OwnerID,
					Currency: account.Currency,
					Balance:  account.Balance,
				}
				store.On("CreateAccount", mock.Anything, arg).
					Return(db.Account{}, mock.Anything)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusInternalServerError",
			body: createAccountRequest{
				OwnerID:  account.OwnerID,
				Currency: account.Currency,
				Balance:  account.Balance,
			},
			build: func(store *mocks.Store) {
				arg := db.CreateAccountParams{
					OwnerID:  account.OwnerID,
					Currency: account.Currency,
					Balance:  account.Balance,
				}
				store.On("CreateAccount", mock.Anything, arg).
					Return(account, sql.ErrConnDone)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	for _, ts := range testCases {
		t.Run(ts.name, func(t *testing.T) {
			store := &mocks.Store{}
			ts.build(store)

			dur, err := time.ParseDuration("1m")
			require.NoError(t, err)

			server, err := NewServer(store, util.Config{
				TokenSymetricKey:    "12345678901234567890123456789012",
				AccessTokenDuration: dur,
			})
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			data, err := json.Marshal(ts.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/account", bytes.NewReader(data))
			require.NoError(t, err)
			req.Header = http.Header{
				"Content-Type": {"application/json"},
			}

			server.router.ServeHTTP(rec, req)
			ts.check(t, rec)
		})
	}
}

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name  string
		url   string
		build func(store *mocks.Store)
		check func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			url:  fmt.Sprintf("/account/%d", account.ID),
			build: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, account.ID).
					Return(account, nil).
					Once()
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchAccount(t, rec.Body, getAccountSuccessResponse{Data: account})
			},
		},
		{
			name: "StatusBadRequestWrongFormatID",
			url:  fmt.Sprintf("/account/%s", "abcde"),
			build: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, account.ID).
					Return(db.Account{}, mock.Anything)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusBadRequestIDCannotBe0",
			url:  fmt.Sprintf("/account/%d", 0),
			build: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, account.ID).
					Return(db.Account{}, mock.Anything)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusNotFound",
			url:  fmt.Sprintf("/account/%d", account.ID),
			build: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, account.ID).
					Return(db.Account{}, sql.ErrNoRows)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			name: "StatusInternalServerError",
			url:  fmt.Sprintf("/account/%d", account.ID),
			build: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, account.ID).
					Return(db.Account{}, sql.ErrConnDone)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	for _, ts := range testCases {
		t.Run(ts.name, func(t *testing.T) {
			store := &mocks.Store{}
			ts.build(store)

			dur, err := time.ParseDuration("1m")
			require.NoError(t, err)

			server, err := NewServer(store, util.Config{
				TokenSymetricKey:    "12345678901234567890123456789012",
				AccessTokenDuration: dur,
			})
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, ts.url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(rec, req)
			ts.check(t, rec)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.GenRandomNum(1, 10000),
		Balance:  util.GenRandomMoney(),
		Currency: util.GenRandomCurrency(),
	}
}

func requireBodyMatchAccount[V createUserSuccessResponse | getAccountErrorResponse | createTransferSuccessResponse | fetchAccountSuccessResponse | createAccountSuccessResponse | getAccountSuccessResponse](t *testing.T, body *bytes.Buffer, res V) {
	bodyData, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotData V
	err = json.Unmarshal(bodyData, &gotData)
	require.NoError(t, err)
	require.Equal(t, res, gotData)
}
