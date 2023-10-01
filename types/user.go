package types

import "time"

type (
	CreateUserRequest struct {
		Username string `json:"username" validate:"required,alpha,min=3,max=32"`
		Password string `json:"password" validate:"required,min=8"`
		FullName string `json:"full_name" validate:"required,min=3,max=32"`
		Email    string `json:"email" validate:"required,email"`
	}

	CreateUserResponse struct {
		Username          string    `json:"username"`
		FullName          string    `json:"full_name"`
		Email             string    `json:"email"`
		PasswordChangedAt time.Time `json:"password_changed_at"`
		CreatedAt         time.Time `json:"created_at"`
	}

	GetUserRequest struct {
		Username string `json:"username" validate:"required,min=3,max=32"`
	}

	LoginUserRequest struct {
		Username string `json:"username" validate:"required,min=3,max=32"`
		Password string `json:"password" validate:"required,min=8"`
	}

	LoginUserResponse struct {
		AccessToken string `json:"access_token"`
		User        CreateUserResponse
	}
)
