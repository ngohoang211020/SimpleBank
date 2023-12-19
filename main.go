package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/gapi"
	"simplebank/pb/user"
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
	runDBMigration(util.Config.MigrationURL, util.Config.DBSource)
	runGrpcServer(store)
}

func runDBMigration(migrationURL string, dbSource string) {
	// New returns a new Migrate instance from a source URL and a database URL.
	// The URL scheme is defined by each driver.
	//Use from file
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to run migrate up:", err)
	}

	log.Println("db migrated successfully")
}

func runGinServer(store db.Store) {
	server, err := api.NewServer(&util.Config, store)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}
	err = server.Start(util.Config.GinPort)
	if err != nil {
		log.Fatal("Cannot connect to server", err)
	}
}

func runGrpcServer(store db.Store) {
	flag.Parse()
	server, err := gapi.NewGrpcServer(&util.Config, store)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	grpcServer := grpc.NewServer()
	user.RegisterAuthServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", util.Config.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("start gRPC server at %s", lis.Addr().String())
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatal("Cannot start gRPC server", err)
	}

}
