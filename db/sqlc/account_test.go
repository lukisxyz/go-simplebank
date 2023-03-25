package db

import (
	"context"
	"testing"
	"time"

	"github.com/flukis/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createDummyAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner: util.GenRandomOwner(),
		Balance: util.GenRandomMoney(),
		Currency: util.GenRandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Balance, args.Balance)
	require.Equal(t, account.Owner, args.Owner)
	require.Equal(t, account.Currency, args.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createDummyAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createDummyAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createDummyAccount(t)

	args := UpdateBalanceAccountParams{
		ID: account1.ID,
		Balance: account1.Balance + 758786,
	}

	account2, err := testQueries.UpdateBalanceAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.NotEqual(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestFetchAccounts(t *testing.T) {
	for i := 0; i< 10; i++ {
		createDummyAccount(t)
	}

	arg := FetchAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.FetchAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}