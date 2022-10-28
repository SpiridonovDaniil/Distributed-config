package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var up = flag.Bool("up", true, "true if up, else down, default true")

func main() {
	//cfg, err := config.LoadConfig("")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//dbCfg := cfg.DataBase

	flag.Parse()

	db, err := sql.Open("postgres",
		fmt.Sprint("host=db port=5432 user=user password=test dbname=config sslmode=disable"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if *up {
		if err := goose.Up(db, "migration"); err != nil {
			log.Fatal(err)
		}
		log.Println("upped")

		return
	}

	if err := goose.Down(db, "migration"); err != nil {
		log.Fatal(err)
	}
	log.Println("downed")
}
