package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/SpiridonovDaniil/Distributed-config/internal/config"
)

const migrationPackage = "migration"

const driverName = "postgres"

var up = flag.Bool("up", true, "true if up, else down, default true")

func main() {
	flag.Parse()

	cfg := config.Read()

	db, err := sql.Open(driverName,
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Postgres.Address,
			cfg.Postgres.Port,
			cfg.Postgres.User,
			cfg.Postgres.Pass,
			cfg.Postgres.Db,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.SetDialect(driverName); err != nil {
		log.Fatal(err)
	}

	if *up {
		if err := goose.Up(db, migrationPackage); err != nil {
			log.Fatal(err)
		}
		log.Println("upped")

		return
	}

	if err := goose.Down(db, migrationPackage); err != nil {
		log.Fatal(err)
	}
	log.Println("downed")
}
