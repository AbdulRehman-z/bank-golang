package api

import (
	"fmt"

	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/gofiber/fiber/v2"
)

// type createAccountRequest struct {
// 	Owner    string `json:"owner"`
// 	Currency string `json:"currency"`
// }

func (server *Server) createAccountHandler(c *fiber.Ctx) error {

	// var account createAccountRequest

	// arg := db.CreateAccountParams{
	// 	Owner:    req.Owner,
	// 	Currency: req.Currency,
	// 	Balance:  0,
	// }

	// account, err := server.store.CreateAccount(c, arg)

	// get the request body
	// var req util.CreateAccountRequest
	req := new(util.CreateAccountRequest)
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Account created successfully",
		"data":    req,
	})
}
