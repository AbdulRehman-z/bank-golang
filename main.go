package main

import (
	"context"
	"database/sql"

	"net"
	"net/http"
	"os"

	"github.com/AbdulRehman-z/bank-golang/api"
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/gapi"
	"github.com/AbdulRehman-z/bank-golang/mail"
	"github.com/AbdulRehman-z/bank-golang/pb"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/AbdulRehman-z/bank-golang/worker"
	"github.com/golang-migrate/migrate/v4"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/AbdulRehman-z/bank-golang/doc/statik"
	_ "github.com/lib/pq"
)

func main() {
	util.ClearConsole()
	// Load env variables
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.ENVIRONMENT == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DB_DRIVER, config.DB_URL)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	store := db.NewStore(conn)
	runDbMigration(config.DB_MIGRATION_URL, config.DB_URL)

	redisClientOpts := asynq.RedisClientOpt{
		Addr: config.REDIS_ADDR,
	}
	taskDistributor := worker.NewTaskDistributor(redisClientOpts)

	mailer := mail.NewGmailSender("Abdul Rehman", config.FROM_EMAIL_ADDRESS, config.APP_PASSWORD)

	go runRedisTaskProcessor(redisClientOpts, store, mailer)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func runRedisTaskProcessor(opts asynq.RedisClientOpt, store db.Store, mailer mail.Mail) {
	taskProcessor := worker.NewRedisTaskProcessor(opts, store, mailer)
	log.Info().Msg("Starting redis task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start redis task processor")
	}
}

func runDbMigration(sourceURL string, dbURL string) {

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create migration")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to migrate db")
	}
}

func runFiberServer(config *util.Config, store db.Store) {
	server, err := api.NewServer(*config, store)
	log.Info().Msg("Starting fiber server on " + config.LISTEN_ADDR)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create fiber server")
	}

	err = server.Start(&config.LISTEN_ADDR)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func runGrpcServer(config *util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(*config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create gRPC server")
	}

	listener, err := net.Listen("tcp", config.GRPC_ADDR)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	// initialize grpc server

	grpcLogger := grpc.UnaryInterceptor(gapi.Logger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterBankServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Info().Msg("Starting grpc server on " + config.GRPC_ADDR)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}

func runGatewayServer(config *util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(*config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create gRPC server")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOptions)
	err = pb.RegisterBankServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register gateway server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create statik file system")
	}
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFS)))

	listener, err := net.Listen("tcp", config.GRPC_GATEWAY_ADDR)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	log.Info().Msg("Starting grpc gateway server on " + config.GRPC_GATEWAY_ADDR)
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
