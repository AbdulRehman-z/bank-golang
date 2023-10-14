package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/AbdulRehman-z/bank-golang/api"
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/gapi"
	"github.com/AbdulRehman-z/bank-golang/pb"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/lib/pq"
)

func main() {
	util.ClearConsole()
	// Load env variables
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load env: ", err)
	}
	conn, err := sql.Open(config.DB_DRIVER, config.DB_URL)
	if err != nil {
		log.Fatal("failed to connect: ", err)
	}
	store := db.NewStore(conn)

	// runFiberServer(config, store)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runFiberServer(config *util.Config, store db.Store) {
	server, err := api.NewServer(*config, store)
	fmt.Println("Starting http server...")
	if err != nil {
		log.Fatal("failed to create http server: ", err)
	}

	log.Fatal(server.Start(&config.LISTEN_ADDR))
}

func runGrpcServer(config *util.Config, store db.Store) {
	server, err := gapi.NewServer(*config, store)
	if err != nil {
		log.Fatal("failed to create server: ", err)
	}

	listener, err := net.Listen("tcp", config.GRPC_ADDR)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	// initialize grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterBankServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	fmt.Printf("Starting grpc server on %s...\n", config.GRPC_ADDR)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("failed to create grpc server: ", err)
	}
}

func runGatewayServer(config *util.Config, store db.Store) {
	server, err := gapi.NewServer(*config, store)
	if err != nil {
		log.Fatal("failed to create server: ", err)
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
		log.Fatal("failed to listen: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.GRPC_GATEWAY_ADDR)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	fmt.Printf("Starting grpc gateway server on %s...\n", config.GRPC_GATEWAY_ADDR)
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("failed to create http gateway server: ", err)
	}
}
