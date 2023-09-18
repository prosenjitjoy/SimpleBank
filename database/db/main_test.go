package db

import (
	"context"
	"log"
	"main/util"
	"math/rand"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	cfg, err := util.LoadConfig("../../.env")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}

func randomAccountID(arr []*Account) int64 {
	n := len(arr)
	return arr[rand.Intn(n)].ID
}
