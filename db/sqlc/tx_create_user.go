package db

import (
	"context"
)

// CreateuserTxParams contains the input parameters of the transfer transaction
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// CreateuserTxResult is the resultant structure for a successful db transaction
type CreateuserTxResult struct {
	User User
}

// var txKey = struct{}{}

// CreateuserTx transfer the amount from one account to other
// It creates a transfer record, add account entries, and update acccount's balance within a single db transaction
func (store *SqlStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateuserTxResult, error) {
	var result CreateuserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err

}
