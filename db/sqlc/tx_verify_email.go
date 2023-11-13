package db

import (
	"context"
	"database/sql"
	"fmt"
)

// VerifyEmailTxParams contains the input parameters of the transfer transaction
type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

// VerifyEmailTxResult is the resultant structure for a successful db transaction
type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

// VerifyEmailTx verifies user email
func (store *SqlStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		fmt.Printf("Email id: %v\n", arg.EmailId)
		fmt.Printf("Secret code: %v\n", arg.SecretCode)

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})
		if err != nil {
			fmt.Printf("update user error: %v\n", err)
			return err
		}

		return err
	})

	return result, err
}
