package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/AbdulRehman-z/bank-golang/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	// Load env variables
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Open a database connection
	testDb, err = sql.Open(config.DB_DRIVER, config.DB_URL)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	defer testDb.Close() // Ensure the database connection is closed when the function exits

	// Initialize testQueries with the database connection
	testQueries = New(testDb)

	// Run all tests
	os.Exit(m.Run())
}
