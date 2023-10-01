package main

import (
	"database/sql"
	"log"

	"github.com/AbdulRehman-z/bank-golang/api"
	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/util"

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
	server, err := api.NewServer(*config, store)
	if err != nil {
		log.Fatal("failed to create server: ", err)
	}
	log.Fatal(server.Start(&config.LISTEN_ADDR))
}
