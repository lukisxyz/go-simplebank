package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/flukis/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	conf, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal(err)
	}
	testDB, err = sql.Open(conf.DBDriver, conf.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
