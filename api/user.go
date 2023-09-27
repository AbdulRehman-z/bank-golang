package api

import (
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/types"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/gofiber/fiber/v2"
)

func (server *Server) createUserHandler(c *fiber.Ctx) error {

	var req types.CreateUserParams
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}

	arg := &db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: req.HashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	// checks if the user already exists
	if _, err := server.store.GetUser(c.Context(), arg.Username); err == nil {
		return fiber.NewError(fiber.StatusBadRequest, USER_ALREADY_EXISTS)
	}

	user, err := server.store.CreateUser(c.Context(), *arg)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, FAILED_TO_CREATE_USER)
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status":  "success",
		"message": "user created",
		"data":    user,
	})
}

func (server *Server) getUserHandler(c *fiber.Ctx) error {

	username := c.Params("username")
	if username == "" {
		return fiber.NewError(fiber.StatusBadRequest, BAD_REQUEST)
	}

	user, err := server.store.GetUser(c.Context(), username)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, FAILED_TO_GET_USER)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "user found",
		"data":    user,
	})
}
