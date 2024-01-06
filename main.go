package main

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/ngohoang211020/simplebank/api"
	db "github.com/ngohoang211020/simplebank/db/sqlc"
	_ "github.com/ngohoang211020/simplebank/doc/statik"
	"github.com/ngohoang211020/simplebank/gapi"
	mail2 "github.com/ngohoang211020/simplebank/mail"
	pb "github.com/ngohoang211020/simplebank/pb/user"
	"github.com/ngohoang211020/simplebank/util"
	"github.com/ngohoang211020/simplebank/worker"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
)

func main() {
	err := util.Config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if util.Config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	connPool, err := pgxpool.New(context.Background(), util.Config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{Addr: util.Config.RedisPort}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	mail := mail2.NewGmailSender(util.Config.EmailSenderName, util.Config.EmailSenderAddress, util.Config.EmailSenderPassword)

	runDBMigration(util.Config.MigrationURL, util.Config.DBSource)

	go runTaskProcessor(redisOpt, store, mail)
	go runGatewayServer(&util.Config, store, taskDistributor, mail)
	runGrpcServer(&util.Config, store, taskDistributor, mail)
}

func runDBMigration(migrationURL string, dbSource string) {
	// New returns a new Migrate instance from a source URL and a database URL.
	// The URL scheme is defined by each driver.
	//Use from file
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runGinServer(store db.Store) {
	server, err := api.NewServer(&util.Config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	err = server.Start(util.Config.GinPort)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mail mail2.EmailSender) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mail)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
func runGrpcServer(config *util.Configuration, store db.Store, distributor worker.TaskDistributor, mail mail2.EmailSender) {
	server, err := gapi.NewGrpcServer(config, store, distributor, mail)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.GrpcPort)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Printf("start gRPC server at %s", lis.Addr().String())
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}

}

func runGatewayServer(config *util.Configuration, store db.Store, taskDistributor worker.TaskDistributor, mail mail2.EmailSender) {
	server, err := gapi.NewGrpcServer(config, store, taskDistributor, mail)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
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
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// fs is subpackge of statik
	statikFs, err := fs.NewWithNamespace("simple_bank")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", util.Config.GinPort)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}
