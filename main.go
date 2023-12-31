package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/ngohoang211020/simplebank/api"
	db "github.com/ngohoang211020/simplebank/db/sqlc"
	_ "github.com/ngohoang211020/simplebank/doc/statik"
	"github.com/ngohoang211020/simplebank/gapi"
	pb "github.com/ngohoang211020/simplebank/pb/user"
	"github.com/ngohoang211020/simplebank/util"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
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
	go runGatewayServer(&util.Config, store)
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
		log.Fatal("Cannot create server:", err)
	}
	err = server.Start(util.Config.GinPort)
	if err != nil {
		log.Fatal("Cannot connect to server:", err)
	}
}

func runGrpcServer(store db.Store) {
	flag.Parse()
	server, err := gapi.NewGrpcServer(&util.Config, store)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", util.Config.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("start gRPC server at %s", lis.Addr().String())
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatal("Cannot start gRPC server:", err)
	}

}

func runGatewayServer(config *util.Configuration, store db.Store) {
	server, err := gapi.NewGrpcServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// fs is subpackge of statik
	statikFs, err := fs.NewWithNamespace("simple_bank")
	if err != nil {
		log.Fatal("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", util.Config.GinPort))
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server: ", err)
	}
}
