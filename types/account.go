package types

type (
	// CreateAccountRequest defines the request body for creating an account
	CreateAccountRequest struct {
		Owner    string `json:"owner" validate:"required,min=3,max=20"` // Required field, min 3 char long max 20
		Currency string `json:"currency" validate:"required"`           // Required field, one of USD, EUR, CAD
	}

	// GetAccountRequest defines the request body for getting an account
	// id must be positive and grater than 0
	GetAccountRequest struct {
		ID int64 `json:"id" validate:"required,gt=0"` // Required field, min 1
	}

	// ListAccountsRequest defines the request params for listing accounts
	ListAccountsRequest struct {
		PageID   int32 `query:"page_id" validate:"required,gte=1"`   // Required field, min 1
		PageSize int32 `query:"page_size" validate:"required,lte=5"` // Required field, min 5
	}

	// UpdateAccountRequest defines the request body for updating an account
	UpdateAccountRequest struct {
		ID      int64 `json:"id" validate:"required,min=1"`      // Required field, min 1
		Balance int64 `json:"balance" validate:"required,min=1"` // Required field, min 1
	}

	// DeleteAccountRequest defines the request body for deleting an account
	DeleteAccountRequest struct {
		ID int64 `json:"id" validate:"required,min=1"` // Required field, min 1
	}
)
