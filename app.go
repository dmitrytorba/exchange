package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var db *sql.DB
var exch *exchange

func main() {
	startDb()
	startExchange()

	// testing creating an order
	order, err := exch.createOrder(1337, 1, 1, "ltc", "sell")
	fmt.Println(order, err)

	startApi()
}

func startExchange() {
	var err error
	exch, err = createExchange()
	if err != nil {
		panic(err)
	}
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
		var err error
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
	}
}
