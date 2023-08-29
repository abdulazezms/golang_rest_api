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

	errs := make(chan error)
	results := make(chan TransferTxResult)

	

	for i := 0; i < n; i++ {
		go func() {
			params := TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID: acc2.ID,
				Amount: transferAmount,
			}
			fmt.Println("Calling transfer with params = ", params)
			result, err := store.TransferTx(context.Background(), params)
			errs <- err 
			results <- result 
		}()
	}

	//check results 
	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
		transferResult := <- results
		fmt.Printf("%+v\n", transferResult)

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
	}
	

}