package api

import (
	"database/sql"
	"fmt"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/token"
	"github.com/AbdulRehman-z/bank-golang/types"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/gofiber/fiber/v2"
)

// createAccountHandler creates a new account
func (server *Server) createTransferHandler(c *fiber.Ctx) error {

	var req types.CreateTransferRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, FAILED_TO_PARSE_BODY)
	}

	// validate the request
	if err := util.CheckValidationErrors(req); err != nil {
		return err
	}

	// NOTE: There are two ways for acheiving the same result. i.e. checking whether the account exists or not
	// 1. using go routines
	// errChannel to catch any possible errors
	// errChannel := make(chan error, 2)
	// // two go routines to fetch result concurrently
	// go func() {
	// 	errChannel <- server.accountExist(c, arg.FromAccountID, req.Currency)
	// }()
	// go func() {
	// 	errChannel <- server.accountExist(c, arg.ToAccountID, req.Currency)
	// }()

	// // any possible errors returnes to errorChannel will be exposed in the following for loop
	// var errors []error
	// for i := 0; i < 2; i++ {
	// 	err := <-errChannel
	// 	if err != nil {
	// 		errors = append(errors, err)
	// 	}
	// }
	// if len(errors) > 0 {
	// 	return errors[0]
	// }

	// 2. using standard way
	fromAccount, err := server.accountExist(c, req.FromAccountID, req.Currency)
	if err != nil {
		return err
	}

	payload := c.Locals(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != payload.Username {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	_, err = server.accountExist(c, req.ToAccountID, req.Currency)
	if err != nil {
		fmt.Println("err: ", err)
		return err
	}

	// create a new account
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// save the account in the database
	tx, err := server.store.TransferTx(c.Context(), arg)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to transfer: %v", err))
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "transfer was successful",
		"data":    tx,
	})
}

// a little helper function that checks whether the account that are being involved in transaction exists or not
func (server *Server) accountExist(c *fiber.Ctx, accountId int64, currency string) (*db.Account, error) {
	account, err := server.store.GetAccount(c.Context(), accountId)
	fmt.Println("======================== No rows")
	fmt.Printf("Account: %v\n", account)
	fmt.Println("======================== No Rows")

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("======================== No rows")
			fmt.Printf("err: %s\n", err)
			fmt.Println("======================== No Rows")
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("account with id %d not found", accountId))
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get account: %v", err))
	}

	fmt.Println("======================== No rows")
	fmt.Printf("%s vs %s : \n", account.Currency, currency)
	fmt.Println("======================== No Rows")

	if account.Currency != currency {
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("account currency mismatch: %s vs %s", account.Currency, currency))
	}

	return &account, nil
}
