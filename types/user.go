package types

type (
	CreateUserRequest struct {
		Username       string `json:"username" validate:"required,min=3,max=32"`
		HashedPassword string `json:"hashed_password" validate:"required,min=8"`
		FullName       string `json:"full_name" validate:"required,min=3,max=32"`
		Email          string `json:"email" validate:"required,email"`
	}

	GetUserRequest struct {
		Username string `json:"username" validate:"required,min=3,max=32"`
	}
)
