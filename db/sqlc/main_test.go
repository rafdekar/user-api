package db

import (
	"database/sql"
	"github.com/rafdekar/user-api/util"
	"log"
	"os"
	"testing"
)

var testDB *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalln("config could be loaded: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("db connection could not be established: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
