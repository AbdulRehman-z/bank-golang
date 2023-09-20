package api

import (
	"fmt"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/types"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/gofiber/fiber/v2"
)

func (server *Server) createAccountHandler(c *fiber.Ctx) error {

	req := new(types.CreateAccountRequest)
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	// validate the request
	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}

	// create a new account
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}
	// save the account in the database
	account, err := server.store.CreateAccount(c.Context(), arg)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Account created successfully",
		"data":    account,
	})

}
