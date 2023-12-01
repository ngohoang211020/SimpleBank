package main

import (
	"database/sql"
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

	conn, err := sql.Open(util.Config.DBDriver, util.Config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(util.Config.Port)
	if err != nil {
		log.Fatal("Cannot connect to server", err)
	}
}
