package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" //an implementation of a driver
	"tutorial.sqlc.dev/app/util"
)


var testQueries *Queries
var testDB *sql.DB
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..", "yml", "config")
	if err != nil {
		log.Fatal("Cannot load configuration: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}