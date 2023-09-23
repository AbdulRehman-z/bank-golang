package types

type (
	CreateTransferRequest struct {
		FromAccountID int64  `json:"from_account_id" validate:"required,min=1"`
		ToAccountID   int64  `json:"to_account_id" validate:"required,min=1"`
		Amount        int64  `json:"amount" validate:"required,gt=0"`
		Currency      string `json:"currency" validate:"required,oneof=USD EUR CAD"`
	}
)
