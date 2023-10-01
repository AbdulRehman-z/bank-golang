package api

import (
	"errors"
	"fmt"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/token"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *fiber.App
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

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

	tokenMaker, err := token.NewPasetoMaker(config.SYMMETRIC_KEY)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}
	server := &Server{
		store:      store,
		router:     app,
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.setupRoutes(app)

	return server, nil

}

func (server *Server) Start(listenAddr *string) error {
	return server.router.Listen(*listenAddr)
}

func (server *Server) setupRoutes(app *fiber.App) {

	auth := app.Group("/v1", AuthMiddleware(server.tokenMaker))

	auth.Post("/accounts", server.createAccountHandler)
	auth.Get("/accounts/:id", server.getAccountHandler)
	auth.Get("/accounts", server.listAccountsHandler)
	auth.Put("/accounts", server.updateAccountHandler)
	auth.Delete("/accounts/:id", server.deleteAccountHandler)
	auth.Post("/transfers", server.createTransferHandler)

	app.Post("/users", server.createUserHandler)
	app.Post("/users/login", server.loginUserHandler)
}
