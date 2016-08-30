package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var db *sql.DB

func main() {
	startDb()
	startExchange()
	startApi()
}

func startExchange() {
	e, err := createExchange()
	if err != nil {
		panic(err)
	}
	fmt.Println(e)
}

func startApi() {
	err := api()
	if err != nil {
		panic(err)
	}
}

func startDb() {
	connStr := os.Getenv("EXCHANGEDB")
	if connStr == "" {
		panic("PSQL environment variable missing (export EXCHANGEDB='postgres://exchange:xNzoA3ZNfTe89Kqp2h@localhost/exchange?sslmode=disable')")
	} else {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
	}
}
