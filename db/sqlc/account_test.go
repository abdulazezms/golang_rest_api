package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"tutorial.sqlc.dev/app/util"
)

func createTestAccount(arg CreateAccountParams) (Account, error) {
	account, err := testQueries.CreateAccount(context.Background(), arg)
	return account, err
}

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}
	account, err := createTestAccount(arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Balance, account.Balance)
	require.NotZero(t, account.ID, 1)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}
	account1, err := createTestAccount(arg)
	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
	require.Equal(t, account1.Currency, account2.Currency)
}

func TestUpdateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}
	account1, err := createTestAccount(arg)
	require.NoError(t, err)

	arg2 := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomAmount(),
		Owner: util.RandomOwner(),
	}
	err = testQueries.UpdateAccount(context.Background(), arg2)
	require.NoError(t, err)

	account1, err = testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.ID, arg2.ID)
	require.Equal(t, account1.Balance, arg2.Balance)
	require.Equal(t, account1.Currency, arg.Currency)
	require.Equal(t, account1.Owner, arg2.Owner)
}

func TestDeleteAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}
	account1, err := createTestAccount(arg)
	require.NoError(t, err)

	err = testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		arg := CreateAccountParams{
			Owner:    util.RandomOwner(),
			Balance:  util.RandomAmount(),
			Currency: util.RandomCurrency(),
		}
		createTestAccount(arg)
	}

	accounts, err := testQueries.ListAccounts(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(accounts), 10)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
