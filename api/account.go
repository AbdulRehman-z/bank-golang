package api

import (
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

type createAccountRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
}

func (server *Server) createAccountHandler(ctx *fiber.Ctx) (*db.Account, error) {

	var req createAccountRequest

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(c, arg)

}
