package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"tutorial.sqlc.dev/app/util"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	n := 5
	transferAmount := util.RandomAmount()
	acc1, err := createTestAccount(CreateAccountParams{
		util.RandomOwner(),
		util.RandomAmount(),
		util.RandomCurrency(),
	})
	require.NoError(t, err)

	acc2, err := createTestAccount(CreateAccountParams{
		util.RandomOwner(),
		util.RandomAmount(),
		util.RandomCurrency(),
	})
	require.NoError(t, err)

	fmt.Println(">> before: ", acc1.Balance, acc2.Balance)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			params := TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        transferAmount,
			}

			result, err := store.TransferTx(context.Background(), params)
			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		transferResult := <-results

		//check transfer
		require.NotEmpty(t, transferResult)
		require.Equal(t, transferResult.Transfer.FromAccountID, acc1.ID)
		require.Equal(t, transferResult.Transfer.ToAccountID, acc2.ID)
		require.Equal(t, transferResult.Transfer.Amount, transferAmount)
		require.NotZero(t, transferResult.Transfer.CreatedAt)
		require.NotZero(t, transferResult.Transfer.ID)

		_, err = store.GetTransfer(context.Background(), transferResult.Transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := transferResult.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -transferAmount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)
		require.NotZero(t, fromEntry.ID)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := transferResult.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, transferAmount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)
		require.NotZero(t, toEntry.ID)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//TODO: CHECK ACCOUNTS' BALANCES
		fromAccount := transferResult.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := transferResult.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)
		fmt.Println(">> Tx: ", fromAccount.Balance, toAccount.Balance)

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%transferAmount == 0) //1*amount, 2*amount, ..., n*amount
		k := int(diff1 / transferAmount)
		require.True(t, 1 <= k && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	//check the final updated balance of the two accounts.

	updatedAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	fmt.Println(">> After: ", updatedAccount1.Balance, updatedAccount2.Balance)

	//check the balances now
	require.Equal(t, acc1.Balance-int64(n)*transferAmount, updatedAccount1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*transferAmount, updatedAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	n := 10
	transferAmount := util.RandomAmount()
	acc1, err := createTestAccount(CreateAccountParams{
		util.RandomOwner(),
		util.RandomAmount(),
		util.RandomCurrency(),
	})
	require.NoError(t, err)

	acc2, err := createTestAccount(CreateAccountParams{
		util.RandomOwner(),
		util.RandomAmount(),
		util.RandomCurrency(),
	})
	require.NoError(t, err)

	fmt.Println(">> before: ", acc1.Balance, acc2.Balance)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		go func(i int) {
			var fromAccountID int64
			var toAccountID int64

			//In case i is even, then make sure to change the IDs to reproduce the deadlock issue.
			if i%2 == 0 {
				fromAccountID = acc1.ID
				toAccountID = acc2.ID
			} else {
				//otherwise, change the IDs to ensure that the deadlock issue is executed again in this scenario.
				fromAccountID = acc2.ID
				toAccountID = acc1.ID
			}
			params := TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        transferAmount,
			}

			_, err := store.TransferTx(context.Background(), params)
			errs <- err
		}(i)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	//check the final updated balance of the two accounts.

	updatedAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	fmt.Println(">> After: ", updatedAccount1.Balance, updatedAccount2.Balance)

	//Since we have made n/2 transactions from account 1 to account 2 and n/2 transactions from account 2 to account 1, then the final
	//balance should be equal to the balance before executing the transactions.
	require.Equal(t, acc1.Balance, updatedAccount1.Balance)
	require.Equal(t, acc2.Balance, updatedAccount2.Balance)
}

