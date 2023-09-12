package db

import (
	"context"
	"time"
)

func createTransfer() Transfer {

	args := CreateTransferParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        10,
		CreatedAt:     time.Now(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
}
