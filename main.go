package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/AbdulRehman-z/bank-golang/api"
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/util"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:password@localhost:5432/bankDb?sslmode=disable"
)

func main() {
	util.ClearConsole()
	listenAddr := flag.String("listenAddr", ":8080", "server listen address")
	flag.Parse()

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("failed to connect: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(&store)
	fmt.Println("Starting server on")
	log.Fatal(server.Start(*listenAddr))
}
