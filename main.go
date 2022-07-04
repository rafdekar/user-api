package main

import (
	"database/sql"
	"github.com/rafdekar/user-api/api"
	db "github.com/rafdekar/user-api/db/sqlc"
	"github.com/rafdekar/user-api/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln("could not load config file: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("db connection could not be established: ", err)
	}

	server := api.NewServer(db.New(conn))

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalln("server could not be started: ", err)
	}
}
