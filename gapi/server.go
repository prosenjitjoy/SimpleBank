package gapi

import (
	"fmt"
	"main/database/db"
	"main/pb"
	"main/token"
	"main/util"
	"main/worker"
)

// Server serves gRPC request for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          *util.ConfigDatabase
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server
func NewServer(store db.Store, cfg *util.ConfigDatabase, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPASETOMaker(cfg.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          cfg,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
