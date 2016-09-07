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
	exch.books["ltc"].insert(createOrder("john", 100, 100, BUY))
	exch.books["ltc"].insert(createOrder("john", 100, 101, SELL))
	startApi()
}

func startExchange() {
	exch = createExchange()
}

func startApi() {
	fmt.Println("API starting...")
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
