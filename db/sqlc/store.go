package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateuserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

type SqlStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func addMoney(ctx context.Context,
	q *Queries,
	account1ID int64,
	amount1 int64,
	account2ID int64,
	amount2 int64) (account1 Account, account2 Account, err error) {

	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     account1ID,
		Amount: strconv.FormatInt(amount1, 10),
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     account2ID,
		Amount: strconv.FormatInt(amount2, 10),
	})

	return
}
