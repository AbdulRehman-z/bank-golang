package api

import (
	"database/sql"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/types"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/gofiber/fiber/v2"
)

func (server *Server) createUserHandler(c *fiber.Ctx) error {
	var req types.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}

	arg := db.CreateUserRequest{
		Username:       req.Username,
		HashedPassword: req.HashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(c.Context(), arg)
	if err != nil {
		if err == sql.ErrConnDone {
			return fiber.NewError(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR)
		}
		return fiber.NewError(fiber.StatusBadRequest, USER_ALREADY_EXISTS)
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status":  "success",
		"message": "user created",
		"data":    user,
	})
}

func (server *Server) getUserHandler(c *fiber.Ctx) error {

	// parse the params
	var req types.GetUserRequest
	if err := c.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	// validate the request
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

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "user found",
		"data":    user,
	})
}
