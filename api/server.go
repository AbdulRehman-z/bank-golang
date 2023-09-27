package api

import (
	"errors"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	store  db.Store
	router *fiber.App
	// httpEngine http.ServeMux // NOTE: ONLY FOR TESTING
}

func NewServer(store db.Store) *Server {

	app := fiber.New(fiber.Config{
		// Global custom error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500 And Error message defaults to "Internal Server Error"
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			// Retrieve the custom status code and message if it's an fiber.*Error
			var e *fiber.Error
			if errors.As(err, &e) {
				message = e.Message
				code = e.Code
			}

			// send the error as json
			return ctx.Status(code).JSON(fiber.Map{
				"success": false,
				"message": message,
			})

			// Return from handler
		},
	})

	app.Use(func(c *fiber.Ctx) error {
		// only json is allowed over post requests
		if c.Method() == "POST" || c.Method() == "PUT" {
			if c.Get("Content-Type") != "application/json" {
				return fiber.NewError(fiber.StatusUnsupportedMediaType, "Content-Type must be application/json")
			}
		}

		return c.Next()
	})

	server := &Server{
		store:  store,
		router: app,
	}

	app.Post("/accounts", server.createAccountHandler)
	app.Get("/accounts/:id", server.getAccountHandler)
	app.Get("/accounts", server.listAccountsHandler)
	app.Put("/accounts", server.updateAccountHandler)
	app.Delete("/accounts/:id", server.deleteAccountHandler)

	app.Post("/transfers", server.createTransferHandler)

	app.Post("/users", server.createUserHandler)
	app.Get("/users/:username", server.getUserHandler)

	return server

}

func (server *Server) Start(listenAddr string) error {
	return server.router.Listen(listenAddr)
}
