package api

import (
	"database/sql"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/types"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func (server *Server) createUserHandler(c *fiber.Ctx) error {
	var req types.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR)
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(c.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				return fiber.NewError(fiber.StatusBadRequest, BAD_REQUEST)
			}
		}
		return fiber.NewError(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR)
	}

	response := types.CreateUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status":  "success",
		"message": "user created",
		"data":    response,
	})
}

func (server *Server) loginUserHandler(c *fiber.Ctx) error {
	var req types.LoginUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}

	user, err := server.store.GetUser(c.Context(), req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.ACCESS_TOKEN_DURATION)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR)
	}

	response := types.LoginUserResponse{
		AccessToken: accessToken,
		User: types.CreateUserResponse{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: user.PasswordChangedAt,
			CreatedAt:         user.CreatedAt,
		},
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "user logged in",
		"data":    response,
	})
}
