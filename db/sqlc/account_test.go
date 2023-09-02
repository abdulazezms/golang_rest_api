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
	argUser := getRandomUserParams()
	createTestUser(argUser)
	arg := CreateAccountParams{
		Owner:    argUser.Username,
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
	argUser := getRandomUserParams()
	createTestUser(argUser)
	arg := CreateAccountParams{
		Owner:    argUser.Username,
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}
	acc, err := createTestAccount(arg)
	require.NoError(t, err)
	accRetrieved, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accRetrieved)
	require.Equal(t, acc.ID, accRetrieved.ID)
	require.Equal(t, acc.Owner, accRetrieved.Owner)
	require.Equal(t, acc.Balance, accRetrieved.Balance)
	require.Equal(t, acc.CreatedAt, accRetrieved.CreatedAt)
	require.Equal(t, acc.Currency, accRetrieved.Currency)
}

func TestUpdateAccount(t *testing.T) {
	argUser := getRandomUserParams()
	createTestUser(argUser)

	arg := CreateAccountParams{
		Owner:    argUser.Username,
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}
	acc, err := createTestAccount(arg)
	require.NoError(t, err)

	arg2 := UpdateAccountParams{
		ID:      acc.ID,
		Balance: util.RandomAmount(),
	}
	acc, err = testQueries.UpdateAccount(context.Background(), arg2)

	require.NoError(t, err)
	require.Equal(t, acc.ID, arg2.ID)
	require.Equal(t, acc.Balance, arg2.Balance)
	require.Equal(t, acc.Currency, arg.Currency)
	require.Equal(t, acc.Owner, arg.Owner)
}

func TestDeleteAccount(t *testing.T) {
	argUser := getRandomUserParams()
	createTestUser(argUser)

	arg := CreateAccountParams{
		Owner:    argUser.Username,
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}
	acc, err := createTestAccount(arg)
	require.NoError(t, err)

	err = testQueries.DeleteAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	accRetrieved, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accRetrieved)

}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		argUser := getRandomUserParams()
		createTestUser(argUser)

		arg := CreateAccountParams{
			Owner:    argUser.Username,
			Balance:  util.RandomAmount(),
			Currency: util.RandomCurrency(),
		}
		account, err := createTestAccount(arg)
		require.NoError(t, err)
		lastAccount = account
	}

	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{Owner: lastAccount.Owner, Limit: 100, Offset: 0})
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, account.Owner, lastAccount.Owner)
	}

}
