package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Yelsnik/trackinginventory/util"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err := sql.Open(config.DBDRIVER, config.DBSOURCE)

	if err != nil {
		log.Fatal("could not connect", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
