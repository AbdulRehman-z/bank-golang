package db

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> Before:", account1.Balance, account2.Balance)

	errChannel := make(chan error)
	resultChannel := make(chan TransferTxResult)

	// Run n concurrent transfer transactions
	n := 2
	amount := int64(10.00)

	for i := 0; i < n; i++ {

		// txName := fmt.Sprintf("tx %d", i+1)

		go func() {
			// ctx := context.WithValue(context.Background(), txKey, txName)

			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errChannel <- err
			resultChannel <- result
		}()
	}

	// Check results
	exited := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errChannel
		require.NoError(t, err)

		result := <-resultChannel
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		transAmount, err := strconv.ParseInt(transfer.Amount, 10, 64)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transAmount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries (to & from)
		toEntry := result.ToEntry
		toAmmount, err := strconv.ParseInt(toEntry.Amount, 10, 64)
		require.NoError(t, err)
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, amount, toAmmount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		fromAmmount, err := strconv.ParseInt(fromEntry.Amount, 10, 64)
		require.NoError(t, err)
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, amount, -fromAmmount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// check accounts' balance
		balance1, err := strconv.ParseFloat(account1.Balance, 64)
		require.NoError(t, err)
		balance2, err := strconv.ParseFloat(fromAccount.Balance, 64)
		require.NoError(t, err)

		balance3, err := strconv.ParseFloat(toAccount.Balance, 64)
		require.NoError(t, err)
		balance4, err := strconv.ParseFloat(account2.Balance, 64)
		require.NoError(t, err)
		diff1 := balance1 - balance2
		diff2 := balance3 - balance4

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, int64(diff1)%int64(amount) == 0)

		k := int(diff1 / float64(amount))
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, exited, k)
		exited[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	acc1Balance, err := strconv.ParseFloat(account1.Balance, 64)
	require.NoError(t, err)
	acc2Balance, err := strconv.ParseFloat(account2.Balance, 64)
	require.NoError(t, err)

	updatedAccount1Balance, err := strconv.ParseFloat(updatedAccount1.Balance, 64)
	require.NoError(t, err)

	updatedAccount2Balance, err := strconv.ParseFloat(updatedAccount2.Balance, 64)
	require.NoError(t, err)

	fmt.Println(">> Compare:", int64(acc1Balance), int64(updatedAccount1Balance))

	require.Equal(t, int64(acc1Balance)-int64(n)*int64(amount), int64(updatedAccount1Balance))
	require.Equal(t, int64(acc2Balance)+int64(n)*int64(amount), int64(updatedAccount2Balance))
	fmt.Println(">> After:", updatedAccount1.Balance, updatedAccount2.Balance)

}

func TestTransferTxDeadLock(t *testing.T) {

	store := NewStore(testDb)

	errorChannel := make(chan error)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(10.00)
	n := 10

	fmt.Println(">> Before:", account1.Balance, account2.Balance)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {

			fmt.Println(">> From:", fromAccountID, "To:", toAccountID)
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errorChannel <- err
		}()

	}

	for i := 0; i < n; i++ {
		errors := <-errorChannel
		require.NoError(t, errors)
	}

	// check if account's balances are equal after transactions
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">> After:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
