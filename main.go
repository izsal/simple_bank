package main

import (
	"context"
	"log"

	"github.com/izsal/simple_bank/api"
	db "github.com/izsal/simple_bank/db/sqlc"
	"github.com/izsal/simple_bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	ctx := context.Background()

	if err != nil {
		log.Fatalf("Unable to parse connString: %v", err)
	}

	pool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	store := db.NewStore(pool)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
