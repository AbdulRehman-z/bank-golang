package api

import (
	"database/sql"
	"fmt"
	"strconv"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/types"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/gofiber/fiber/v2"
)

// Error messages
const (
	FAILED_TO_PARSE_BODY         = "failed to parse body"
	FAILED_TO_PARSE_QUERY_PARAMS = "failed to parse query parameters"
	FAILED_TO_DELETE_ACCOUNT     = "failed to delete account"
	FAILED_TO_UPDATE_ACCOUNT     = "failed to update account"
	FAILED_TO_CREATE_ACCOUNT     = "failed to create account"
	FAILED_TO_GET_ACCOUNT        = "failed to get account"
	FAILED_TO_LIST_ACCOUNTS      = "failed to list accounts"
	ACCOUNT_NOT_FOUND            = "account not found"
	INTERNAL_SERVER_ERROR        = "internal server error"
	CURRENCY_MISMATCH            = "currency mismatched"
	FAILED_TO_CREATE_USER        = "failed to create user"
	FAILED_TO_GET_USER           = "failed to get user"
	USER_ALREADY_EXISTS          = "user already exists"
	BAD_REQUEST                  = "bad request"
)

// createAccountHandler creates a new account
func (server *Server) createAccountHandler(c *fiber.Ctx) error {

	var req types.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println("err: ", err)
		return fiber.NewError(fiber.StatusBadRequest, FAILED_TO_PARSE_BODY)
	}

	// validate the request
	if err := util.CheckValidationErrors(req); err != nil {
		fmt.Println("err: ", err)
		return err
	}

	// create a new account
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  "0",
		Currency: req.Currency,
	}

	// save the account in the database
	account, err := server.store.CreateAccount(c.Context(), arg)
	if err != nil {
		fmt.Println("err: ", err)
		return fiber.NewError(fiber.StatusBadRequest, "account with this owner already exists")
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
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
			return fiber.NewError(fiber.StatusNotFound, ACCOUNT_NOT_FOUND)
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR)
		}
	}

	return c.JSON(&fiber.Map{
		"success": true,
		"message": "Account retrieved successfully",
		"data":    account,
	})

}

// listAccountsHandler lists all accounts
func (server *Server) listAccountsHandler(c *fiber.Ctx) error {
	var query types.ListAccountsRequest

	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, FAILED_TO_PARSE_QUERY_PARAMS)
	}

	if query.PageID == 0 && query.PageSize == 0 {
		return fiber.NewError(fiber.StatusBadRequest, FAILED_TO_PARSE_QUERY_PARAMS)
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
			return fiber.NewError(fiber.StatusNotFound, ACCOUNT_NOT_FOUND)
		}
		return fiber.NewError(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "Accounts retrieved successfully",
		"data":    accounts,
	})
}

// updateAccountHandler updates an account
func (server *Server) updateAccountHandler(c *fiber.Ctx) error {

	var req types.UpdateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, FAILED_TO_PARSE_BODY)
	}

	// validate the request
	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}

	// update the account in the database
	// convert req.Balance to string

	account, err := server.store.UpdateAccount(c.Context(), db.UpdateAccountParams{
		ID:      req.ID,
		Balance: strconv.FormatInt(req.Balance, 10),
	})

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
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
		return fiber.NewError(fiber.StatusInternalServerError, FAILED_TO_DELETE_ACCOUNT)
	}

	return c.JSON(&fiber.Map{
		"success": true,
		"message": "Account deleted successfully",
	})

}
