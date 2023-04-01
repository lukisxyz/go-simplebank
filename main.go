package main

import (
	"database/sql"
	"log"

	"github.com/flukis/simplebank/api"
	db "github.com/flukis/simplebank/db/sqlc"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)


const (
	addr = "0.0.0.0:3000"
	dbDriver = "postgres"
	dbSource = "postgresql://root:wap12345@127.0.0.1:5432/simplebank?sslmode=disable"
)

func main() {
	dbConn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	v := validator.New()

	store := db.NewStore(dbConn)
	server := api.NewServer(store, v)
	server.Start(addr)
}