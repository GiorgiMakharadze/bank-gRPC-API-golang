package main

import (
	"database/sql"
	"log"

	"github.com/GiorgiMakharadze/bank-API-golang/api"
	db "github.com/GiorgiMakharadze/bank-API-golang/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/bank?sslmode=disable"
	serverAddress = "localhost:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect database", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
