package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"tutorial.sqlc.dev/app/api"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/util"
)

func main() {
	config, err := util.LoadConfig(".", "yml", "config")
	if err != nil {
		log.Fatal("Cannot load configuration: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("error while connecting to DB: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("error while creating new server: ", err)
	}

	err = server.Start(config.ServerAaddress)
}
