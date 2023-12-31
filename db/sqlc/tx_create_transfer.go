package db

import (
	"context"
	"fmt"
	"strconv"
)

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the resultant structure for a successful db transaction
type TransferTxResult struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// var txKey = struct{}{}

// TransferTx transfer the amount from one account to other
// It creates a transfer record, add account entries, and update acccount's balance within a single db transaction
func (store *SqlStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// txName := ctx.Value(txKey)
		// Transfer record
		// fmt.Println(txName, "create transfer 1")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        strconv.FormatInt(arg.Amount, 10),
		})
		if err != nil {
			return fmt.Errorf("context: Transfer, err: %v", err)
		}

		// Add account entries
		// fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    strconv.FormatInt(-arg.Amount, 10),
		})
		if err != nil {
			return fmt.Errorf("context: FromEntry, err: %v", err)
		}

		// fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    strconv.FormatInt(+arg.Amount, 10),
		})
		if err != nil {
			return fmt.Errorf("context: ToEntry, err: %v", err)
		}

		//TODO: Update accounts
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, +arg.Amount)
			if err != nil {
				return fmt.Errorf("context: UpdateAccount, err: %v", err)
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, +arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return fmt.Errorf("context: UpdateAccount, err: %v", err)
			}
		}

		return nil
	})

	return result, err

}
