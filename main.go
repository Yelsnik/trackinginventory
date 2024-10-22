package main

import (
	"database/sql"
	"log"

	"github.com/Yelsnik/trackinginventory/api"
	db "github.com/Yelsnik/trackinginventory/db/sqlc"
	"github.com/Yelsnik/trackinginventory/util"
	_ "github.com/lib/pq"
)

func main() {

	// load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// connect to database
	conn, err := sql.Open(config.DBDRIVER, config.DBSOURCE)
	if err != nil {
		log.Fatal("could not connect", err)
	}

	// create a new server
	db := db.NewStore(conn)
	server, err := api.NewServer(config, db)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// start the server
	err = server.Start(config.PORT)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
