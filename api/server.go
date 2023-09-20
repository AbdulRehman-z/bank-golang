package api

import (
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	store  *db.Store
	router *fiber.App
}

func NewServer(store *db.Store) *Server {
	app := fiber.New()

	return &Server{
		store:  store,
		router: app,
	}
}

func (server *Server) Start(listenAddr string) error {
	return server.router.Listen(listenAddr)
}
