package db

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

const DATABASE_URL = "postgres://postgres:postgres@localhost:5432/bankdb"

func TestMain(m *testing.M) {
	var err error
	testDB, err = pgxpool.New(context.Background(), DATABASE_URL)
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
