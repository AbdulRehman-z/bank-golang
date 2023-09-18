package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, errRollback)
		}
		return err
	}

	return tx.Commit()
}

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
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// txName := ctx.Value(txKey)
		// Transfer record
		// fmt.Println(txName, "create transfer 1")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("context: Transfer, err: %v", err)
		}

		// Add account entries
		// fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("context: FromEntry, err: %v", err)
		}

		// fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    +arg.Amount,
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

func addMoney(ctx context.Context,
	q *Queries,
	account1ID int64,
	amount1 int64,
	account2ID int64,
	amount2 int64) (account1 Account, account2 Account, err error) {

	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     account1ID,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     account2ID,
		Amount: amount2,
	})

	return
}
