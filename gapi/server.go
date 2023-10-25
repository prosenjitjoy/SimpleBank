package gapi

import (
	"fmt"
	"main/database/db"
	"main/pb"
	"main/token"
	"main/util"
)

// Server serves gRPC request for our banking service.
type Server struct {
	config     *util.ConfigDatabase
	store      db.Store
	tokenMaker token.Maker
	pb.UnimplementedSimpleBankServer
}

// NewServer creates a new gRPC server
func NewServer(store db.Store, cfg *util.ConfigDatabase) (*Server, error) {
	tokenMaker, err := token.NewPASETOMaker(cfg.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     cfg,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
