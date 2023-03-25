package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createDummyAccount(t)
	account2 := createDummyAccount(t)

	// concurrent transfer transaction
	n := 5
	amount := int64(10)

	errChan := make(chan error)
	resChan := make(chan TransferTxResult)

	for i := 0; i<n ; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			errChan<-err
			resChan<-result
		}()
	}

	for i := 0; i<n ; i++ {
		err := <- errChan
		require.NoError(t, err)

		res := <- resChan
		require.NotEmpty(t, res)

		// test transfer
		tf := res.Transfer
		require.NotEmpty(t, tf)
		require.Equal(t, account1.ID, tf.FromAccountID)
		require.Equal(t, account2.ID, tf.ToAccountID)
		require.Equal(t, amount, tf.Amount)
		require.NotZero(t, tf.ID)
		require.NotZero(t, tf.CreatedAt)
		_, err = store.GetTransfer(context.Background(), tf.ID)
		require.NoError(t, err)

		// test entry from
		fe := res.FromEntry
		require.NotEmpty(t, fe)
		require.Equal(t, account1.ID, fe.AccountID)
		require.Equal(t, -amount, fe.Amount)
		require.NotZero(t, fe.ID)
		require.NotZero(t, fe.CreatedAt)
		_, err = store.GetEntry(context.Background(), fe.ID)
		require.NoError(t, err)

		// test entry from
		te := res.ToEntry
		require.NotEmpty(t, te)
		require.Equal(t, account2.ID, te.AccountID)
		require.Equal(t, amount, te.Amount)
		require.NotZero(t, te.ID)
		require.NotZero(t, te.CreatedAt)
		_, err = store.GetEntry(context.Background(), te.ID)
		require.NoError(t, err)

		// TODO: test balance
	}
}