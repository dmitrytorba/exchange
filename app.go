package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/redis.v4"
	"os"
)

var db *sql.DB
var exch *exchange
var rd *redis.Client

func main() {
	startDb()
	startRedis()
	startExchange()

	// make sure we get all those orders stored in the database
	exch.loadFromDB()

	// bitfinex.go
	connectBitfinex()

	startApi()
}

func startRedis() {
	rd = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rd.Ping().Result()
	if err != nil {
		panic(err)
	}
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
