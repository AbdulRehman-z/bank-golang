package api

import (
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	store  *db.Store
	router *fiber.App
}

func NewServer(store *db.Store) *Server {

	app := fiber.New(fiber.Config{
		// Global custom error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(util.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})
	server := &Server{
		store:  store,
		router: app,
	}

	app.Post("/accounts", server.createAccountHandler)

	return server
	// return &Server{
	// 	store:  store,
	// 	router: app,
	// }
}

func (server *Server) Start(listenAddr string) error {
	return server.router.Listen(listenAddr)
}
