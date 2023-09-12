package main

import (
	"context"
	"log"
	"main/api"
	"main/database/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

const ()

func main() {
	conn, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(SERVER_ADDR)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
