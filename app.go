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
	//startDb()
	//startExchange()
	//startApi()

	orderbook := createOrderbook()
	orderbook.match(createOrder("jacob", 100, 99, SELL))
	execs := orderbook.match(createOrder("jacob", 200, 100, BUY))
	for i := 0; i < len(execs); i++ {
		exec := execs[i]
		fmt.Printf("executed %v units at %v price from %v\n", exec.amount, exec.price, exec.name)
	}
	orderbook.print()
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
