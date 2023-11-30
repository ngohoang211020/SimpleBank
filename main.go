package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:123456@localhost:5432/simple_bank?sslmode=disable"
	port     = "8000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(port)
	if err != nil {
		log.Fatal("Cannot connect to server", err)
	}
}
