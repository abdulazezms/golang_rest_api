package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"tutorial.sqlc.dev/app/api"
	db "tutorial.sqlc.dev/app/db/sqlc"
)

const (
	dbDriver       = "postgres"
	dbSource       = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAaddress = "0.0.0.0:8080"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("error while connecting to DB: ", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAaddress)
}