package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"
)

func main() {
	err := util.Config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), util.Config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store)

	err = server.Start(util.Config.Port)
	if err != nil {
		log.Fatal("Cannot connect to server", err)
	}
}
