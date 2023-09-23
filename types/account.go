package types

type (
	// CreateAccountRequest defines the request body for creating an account
	CreateAccountRequest struct {
		Owner    string `json:"owner" validate:"required,min=5,max=20"` // Required field, min 5 char long max 20
		Currency string `json:"currency" validate:"required"`           // Required field, one of USD, EUR, CAD
	}

	// GetAccountRequest defines the request body for getting an account
	// id must be positive and grater than 0
	GetAccountRequest struct {
		ID int64 `validate:"required,gt=0"` // Required field, min 1
	}

	// ListAccountsRequest defines the request body for listing accounts
	ListAccountsRequest struct {
		PageID   int32 `validate:"required,gte=1"` // Required field, min 1
		PageSize int32 `validate:"required,lte=5"` // Required field, min 5
	}

	// UpdateAccountRequest defines the request body for updating an account
	UpdateAccountRequest struct {
		ID      int64 `validate:"required,min=1"` // Required field, min 1
		Balance int64 `validate:"required,min=1"` // Required field, min 1
	}

	// DeleteAccountRequest defines the request body for deleting an account
	DeleteAccountRequest struct {
		ID int64 `validate:"required,min=1"` // Required field, min 1
	}
)
