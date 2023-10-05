package main

import (
	"context"
	"log"
	"main/api"
	"main/database/db"
	"main/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := util.LoadConfig(".env")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, cfg)
	if err != nil {
		log.Fatal("cannot initialize server:", err)
	}

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
