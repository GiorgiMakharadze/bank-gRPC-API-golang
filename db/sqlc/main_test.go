package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/GiorgiMakharadze/bank-API-golang/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect database", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
