package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/flukis/simplebank/db/mock"
	db "github.com/flukis/simplebank/db/sqlc"
	"github.com/flukis/simplebank/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTransferAPI(t *testing.T) {

	fromAcc := db.Account{
		ID:       util.GenRandomNum(1, 10000),
		Balance:  util.GenRandomMoney(),
		Currency: "IDR",
	}

	toAcc := db.Account{
		ID:       util.GenRandomNum(1, 10000),
		Balance:  util.GenRandomMoney(),
		Currency: "IDR",
	}

	transfer := generateTransferResult(fromAcc, toAcc, 100)

	type wrongCreateTransferParams struct {
		FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
		ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
		Currency      string `json:"currency" binding:"required,oneof=USD EUR IDR"`
		Amount        string `json:"amount" binding:"requied,gt=0"`
	}

	testCases := []struct {
		name  string
		body  any
		build func(store *mocks.Store)
		check func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			body: createTransferRequest{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Currency:      "IDR",
				Amount:        100,
			},
			build: func(store *mocks.Store) {
				arg := db.TransferTxParams{
					FromAccountID: fromAcc.ID,
					ToAccountID:   toAcc.ID,
					Amount:        100,
				}
				store.On("GetAccount", mock.Anything, fromAcc.ID).
					Return(fromAcc, nil).
					Once()
				store.On("GetAccount", mock.Anything, toAcc.ID).
					Return(toAcc, nil).
					Once()
				store.On("TransferTx", mock.Anything, arg).
					Return(transfer, nil).
					Once()
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchAccount(t, rec.Body, createTransferSuccessResponse{Data: transfer})
			},
		},
		{
			name: "StatusBadRequestValidationError",
			body: createTransferRequest{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Currency:      "IBM",
				Amount:        100,
			},
			build: func(store *mocks.Store) {
				arg := db.TransferTxParams{
					FromAccountID: fromAcc.ID,
					ToAccountID:   toAcc.ID,
					Amount:        100,
				}
				store.On("GetAccount", mock.Anything, fromAcc.ID).
					Return(fromAcc, nil)
				store.On("GetAccount", mock.Anything, toAcc.ID).
					Return(toAcc, nil)
				store.On("TransferTx", mock.Anything, arg).
					Return(transfer, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusBadRequestWrongParam",
			body: wrongCreateTransferParams{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Currency:      "IBM",
				Amount:        "USD",
			},
			build: func(store *mocks.Store) {
				arg := db.TransferTxParams{
					FromAccountID: fromAcc.ID,
					ToAccountID:   toAcc.ID,
					Amount:        100,
				}
				store.On("GetAccount", mock.Anything, fromAcc.ID).
					Return(fromAcc, nil)
				store.On("GetAccount", mock.Anything, toAcc.ID).
					Return(toAcc, nil)
				store.On("TransferTx", mock.Anything, arg).
					Return(transfer, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusTransferInternalServerError",
			body: createTransferRequest{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Currency:      "IDR",
				Amount:        100,
			},
			build: func(store *mocks.Store) {
				arg := db.TransferTxParams{
					FromAccountID: fromAcc.ID,
					ToAccountID:   toAcc.ID,
					Amount:        100,
				}
				store.On("GetAccount", mock.Anything, fromAcc.ID).
					Return(fromAcc, nil).
					Once()
				store.On("GetAccount", mock.Anything, toAcc.ID).
					Return(toAcc, nil).
					Once()
				store.On("TransferTx", mock.Anything, arg).
					Return(transfer, sql.ErrConnDone)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "StatusOKButCurrencyNotSame",
			body: createTransferRequest{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Currency:      "USD",
				Amount:        100,
			},
			build: func(store *mocks.Store) {
				arg := db.TransferTxParams{
					FromAccountID: fromAcc.ID,
					ToAccountID:   toAcc.ID,
					Amount:        100,
				}
				store.On("GetAccount", mock.Anything, fromAcc.ID).
					Return(fromAcc, nil).
					Once()
				store.On("GetAccount", mock.Anything, toAcc.ID).
					Return(toAcc, nil).
					Once()
				store.On("TransferTx", mock.Anything, arg).
					Return(transfer, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "StatusAccountNotFound",
			body: createTransferRequest{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Currency:      "IDR",
				Amount:        100,
			},
			build: func(store *mocks.Store) {
				arg := db.TransferTxParams{
					FromAccountID: fromAcc.ID,
					ToAccountID:   toAcc.ID,
					Amount:        100,
				}
				store.On("GetAccount", mock.Anything, toAcc.ID).
					Return(fromAcc, sql.ErrNoRows).
					Once()
				store.On("GetAccount", mock.Anything, fromAcc.ID).
					Return(toAcc, sql.ErrNoRows).
					Once()
				store.On("TransferTx", mock.Anything, arg).
					Return(transfer, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			name: "StatusAccountInternalServerError",
			body: createTransferRequest{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Currency:      "IDR",
				Amount:        100,
			},
			build: func(store *mocks.Store) {
				arg := db.TransferTxParams{
					FromAccountID: fromAcc.ID,
					ToAccountID:   toAcc.ID,
					Amount:        100,
				}
				store.On("GetAccount", mock.Anything, fromAcc.ID).
					Return(fromAcc, sql.ErrConnDone).
					Once()
				store.On("GetAccount", mock.Anything, toAcc.ID).
					Return(toAcc, nil).
					Once()
				store.On("TransferTx", mock.Anything, arg).
					Return(transfer, nil)
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

			conf := util.Config{}

			server, err := NewServer(store, conf)
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			data, err := json.Marshal(ts.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(data))
			require.NoError(t, err)
			req.Header = http.Header{
				"Content-Type": {"application/json"},
			}

			server.router.ServeHTTP(rec, req)
			ts.check(t, rec)
		})
	}
}

func generateTransferResult(a, b db.Account, amount int64) db.TransferTxResult {
	transfer := db.Transfer{
		ID:            util.GenRandomNum(0, 1000),
		FromAccountID: a.ID,
		ToAccountID:   b.ID,
		Amount:        amount,
	}

	fromEntry := db.Entry{
		ID:        util.GenRandomNum(0, 1000),
		AccountID: a.ID,
		Amount:    -amount,
	}

	toEntry := db.Entry{
		ID:        util.GenRandomNum(0, 1000),
		AccountID: b.ID,
		Amount:    amount,
	}

	a.Balance = a.Balance - amount
	b.Balance = b.Balance + amount

	return db.TransferTxResult{
		Transfer:    transfer,
		FromAccount: a,
		ToAccount:   b,
		FromEntry:   fromEntry,
		ToEntry:     toEntry,
	}
}
