package main

import (
	"database/sql"
	"log"

	"github.com/GiorgiMakharadze/bank-API-golang/api"
	db "github.com/GiorgiMakharadze/bank-API-golang/db/sqlc"
	"github.com/GiorgiMakharadze/bank-API-golang/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect database", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.SERVER_ADDRESS)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
