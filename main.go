package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/AbdulRehman-z/bank-golang/api"
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	util.ClearConsole()
	listenAddr := flag.String("listenaddr", ":8080", "server listen address")
	flag.Parse()

	// Load env variables
	godotenv.Load(".env")
	dbDriver, exists := os.LookupEnv("DB_DRIVER")
	if !exists {
		log.Fatal("DB_DRIVER environment variable not set")
	}
	dbSource, exists := os.LookupEnv("DB_URL")
	if !exists {
		log.Fatal("DB_URL environment variable not set")
	}

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("failed to connect: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	// fmt.Println("Starting server on")
	log.Fatal(server.Start(*listenAddr))
}
