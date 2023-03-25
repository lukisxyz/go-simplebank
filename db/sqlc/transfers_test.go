package db

import (
	"context"
	"testing"
	"time"

	"github.com/flukis/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createDummyTransfer(t *testing.T, accountFrom, accountTo Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID: accountTo.ID,
		Amount: util.GenRandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createDummyAccount(t)
	account2 := createDummyAccount(t)
	createDummyTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createDummyAccount(t)
	account2 := createDummyAccount(t)
	transfer1 := createDummyTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestFetchTransfers(t *testing.T) {
	account1 := createDummyAccount(t)
	account2 := createDummyAccount(t)
	for i := 0; i < 10; i++ {
		createDummyTransfer(t, account1, account2)
	}

	arg := FetchTransferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.FetchTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account2.ID)
	}
}