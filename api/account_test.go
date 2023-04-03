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

	mockdb "github.com/flukis/simplebank/db/mock"
	db "github.com/flukis/simplebank/db/sqlc"
	"github.com/flukis/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name  string
		url   string
		build func(store *mockdb.MockStore)
		check func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			url:  fmt.Sprintf("/accounts/%d", account.ID),
			build: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchAccount(t, rec.Body, getAccountSuccessResponse{Data: account})
			},
		},
		{
			name: "StatusNotFound",
			url:  fmt.Sprintf("/accounts/%d", account.ID),
			build: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			name: "StatusBadRequest",
			url:  fmt.Sprintf("/accounts/%d", 0),
			build: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusBadRequest",
			url:  fmt.Sprintf("/accounts/%s", "abcde"),
			build: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusInternalServerError",
			url:  fmt.Sprintf("/accounts/%d", account.ID),
			build: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			store := mockdb.NewMockStore(ctr)
			tc.build(store)

			// start test server
			server := NewServer(store)
			rec := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(rec, req)
			tc.check(t, rec)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.GenRandomNum(1, 1000),
		Owner:    util.GenRandomOwner(),
		Balance:  util.GenRandomMoney(),
		Currency: util.GenRandomCurrency(),
	}
}

func TestCreateAccount(t *testing.T) {
	account := randomAccount()

	type WrongTypeCreateAccount struct {
		Owner    string
		Balance  string
		Currency string
	}

	testCases := []struct {
		name  string
		body  any
		build func(store *mockdb.MockStore)
		check func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			body: createAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
				Balance:  account.Balance,
			},
			build: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  account.Balance,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchAccount(t, rec.Body, createAccountSuccessResponse{Data: account})
			},
		},
		{
			name: "StatusInternalServerError",
			body: createAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
				Balance:  account.Balance,
			},
			build: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(account, sql.ErrConnDone)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "StatusBadRequest",
			body: createAccountRequest{
				Owner:    account.Owner,
				Currency: "IDR",
				Balance:  account.Balance,
			},
			build: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: "IDR",
					Balance:  account.Balance,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusBadRequest",
			body: WrongTypeCreateAccount{
				Owner:    account.Owner,
				Currency: "IDR",
				Balance:  "IDR",
			},
			build: func(store *mockdb.MockStore) {
				arg := WrongTypeCreateAccount{
					Owner:    account.Owner,
					Currency: "IDR",
					Balance:  "IDR",
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			store := mockdb.NewMockStore(ctr)
			tc.build(store)

			// start test server
			server := NewServer(store)
			rec := httptest.NewRecorder()

			d, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(d))
			require.NoError(t, err)
			req.Header = http.Header{
				"Content-Type": {"application/json"},
			}

			server.router.ServeHTTP(rec, req)
			tc.check(t, rec)
		})
	}
}

func requireBodyMatchAccount[V getAccountErrorResponse | createAccountSuccessResponse | getAccountSuccessResponse](t *testing.T, body *bytes.Buffer, res V) {
	bodyData, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotData V
	err = json.Unmarshal(bodyData, &gotData)
	require.NoError(t, err)
	require.Equal(t, res, gotData)
}
