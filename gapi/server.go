package gapi

import (
	"fmt"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/pb"
	"github.com/AbdulRehman-z/bank-golang/token"
	"github.com/AbdulRehman-z/bank-golang/util"
)

// Server serves gRPC requests
type Server struct {
	pb.UnimplementedBankServiceServer
	tokenMaker token.Maker
	config     util.Config
	store      db.Store
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	token, err := token.NewPasetoMaker(config.SYMMETRIC_KEY)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &Server{
		tokenMaker: token,
		config:     config,
		store:      store,
	}, nil
}
