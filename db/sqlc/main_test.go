package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	// Load env variables
	godotenv.Load("../../.env")
	dbDriver, exists := os.LookupEnv("DB_DRIVER")
	if !exists {
		log.Fatal("DB_DRIVER environment variable not set")
	}
	dbSource, exists := os.LookupEnv("DB_URL")
	if !exists {
		log.Fatal("DB_URL environment variable not set")
	}

	// Open a database connection
	var err error
	testDb, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	defer testDb.Close() // Ensure the database connection is closed when the function exits

	// Initialize testQueries with the database connection
	testQueries = New(testDb)

	// Run all tests
	os.Exit(m.Run())
}
