package gapi

import (
	"fmt"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/pb"
	"github.com/AbdulRehman-z/bank-golang/token"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/AbdulRehman-z/bank-golang/worker"
)

// Server serves gRPC requests
type Server struct {
	pb.UnimplementedBankServiceServer
	tokenMaker      token.Maker
	config          util.Config
	store           db.Store
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {

	token, err := token.NewPasetoMaker(config.SYMMETRIC_KEY)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &Server{
		tokenMaker:      token,
		config:          config,
		store:           store,
		taskDistributor: taskDistributor,
	}, nil
}
