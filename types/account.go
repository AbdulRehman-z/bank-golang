package types

type CreateAccountRequest struct {
	Owner    string `validate:"required,min=5,max=20"`      // Required field, min 5 char long max 20
	Currency string `validate:"required,oneof=USD EUR CAD"` // Required field, and client needs to implement our 'teener' tag format which we'll see later
}
