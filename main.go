package main

import (
	"database/sql"
	"log"

	"github.com/flukis/simplebank/api"
	db "github.com/flukis/simplebank/db/sqlc"
	"github.com/flukis/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	dbConn, err := sql.Open(conf.DBDriver, conf.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(dbConn)
	server, err := api.NewServer(store, conf)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}
	server.Start(conf.ServerAddr)
}
