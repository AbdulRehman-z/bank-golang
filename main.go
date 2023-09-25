package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// long pooling after every 24 hours to update exchange rates, save them in db
	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ticker.C:
			// save new exchange rates in db
			url := "https:openexchangerates.org/api/latest.json?app_id=a7784caacbb24a7b9e6129733000733a"
			// get exchange rates from api
			request, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				panic(err)
			}
			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				panic(err)
			}
			defer response.Body.Close()
			// save exchange rates in db
			fmt.Printf("response status: %v\n", response.Body)
		}
	}

}
