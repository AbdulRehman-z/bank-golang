package api

import (
	"database/sql"
	"fmt"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/types"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/gofiber/fiber/v2"
)

// createAccountHandler creates a new account
func (server *Server) createAccountHandler(c *fiber.Ctx) error {

	req := new(types.CreateAccountRequest)
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to parse body")
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
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Account created successfully",
		"data":    account,
	})

}

// getAccountHandler gets an account by id
func (server *Server) getAccountHandler(c *fiber.Ctx) error {

	req := new(types.GetAccountRequest)
	// get the uri param
	id := c.Params("id")
	req.ID = int64(util.StringToInt(id))

	// validate the request
	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}
	// get the account from the database
	account, err := server.store.GetAccount(c.Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "Account not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Account retrieved successfully",
		"data":    account,
	})

}

// listAccountsHandler lists all accounts
func (server *Server) listAccountsHandler(c *fiber.Ctx) error {
	// Parse the query parameters using c.QueryParser()
	fmt.Println("1st ---------------------------------------------------")
	query := new(types.ListAccountsRequest)
	if err := c.QueryParser(&query); err != nil {
		fmt.Println("2nd ---------------------------------------------------")

		return fiber.NewError(fiber.StatusBadRequest, "failed to parse query parameters")
	}

	// Validate the request
	if err := util.CheckValidationErrors(query); err != nil {
		return err
	}

	// List all accounts from the database
	accounts, err := server.store.ListAccounts(c.Context(), db.ListAccountsParams{
		Limit:  int32(query.PageSize),
		Offset: (int32(query.PageID) - 1) * int32(query.PageSize),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			fiber.NewError(fiber.StatusNotFound, "Account not found")
		} else {
			fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Accounts retrieved successfully",
		"data":    accounts,
	})
}

// updateAccountHandler updates an account
func (server *Server) updateAccountHandler(c *fiber.Ctx) error {

	req := new(types.UpdateAccountRequest)
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	// validate the request
	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}

	// update the account in the database
	account, err := server.store.UpdateAccount(c.Context(), db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	})
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Account updated successfully",
		"data":    account,
	})

}

// deleteAccountHandler deletes an account
func (server *Server) deleteAccountHandler(c *fiber.Ctx) error {

	req := new(types.DeleteAccountRequest)

	// get the uri param
	id := c.Params("id")
	req.ID = int64(util.StringToInt(id))

	// validate the request
	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}

	// delete the account from the database
	err := server.store.DeleteAccount(c.Context(), req.ID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Account deleted successfully",
	})

}
