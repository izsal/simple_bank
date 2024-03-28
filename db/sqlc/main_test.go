package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/izsal/simple_bank/util"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *pgx.Conn

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
