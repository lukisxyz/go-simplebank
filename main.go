package main

import (
	"database/sql"
	"log"

	"github.com/flukis/simplebank/api"
	db "github.com/flukis/simplebank/db/sqlc"
	"github.com/flukis/simplebank/util"
	"github.com/go-playground/validator/v10"
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

	v := validator.New()

	store := db.NewStore(dbConn)
	server := api.NewServer(store, v)
	server.Start(conf.ServerAddr)
}
