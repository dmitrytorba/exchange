package main

import (
	"math/rand"
)

type exchange struct {
	books      map[string]*orderbook
	currencies []string
}

func createExchange() *exchange {
	e := &exchange{
		books:      make(map[string]*orderbook),
		currencies: make([]string, 0, 10),
	}

	// get our currency types from the database
	e.currencies = make([]string, 0, 10)
	rows, err := db.Query("SELECT e.enumlabel FROM pg_enum e JOIN pg_type t ON e.enumtypid = t.oid WHERE t.typname = 'currency'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			panic(err)
		}

		e.currencies = append(e.currencies, name)
	}

	// cycle through the config currencies, creating an orderbook for each
	for i := 0; i < len(e.currencies); i++ {
		currency := e.currencies[i]
		e.books[currency] = createOrderbook()
		go e.books[currency].readTicker()
	}

	return e
}

func (e *exchange) loadFromDB() {
	// fill the orderbook with persisted orders
	orders, err := getAllOrders()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(orders); i++ {

		// I added the following if statement to stop a weird situation:
		// lets imagine that we have specified two currencies (btc, ltc)
		// but for some reason the database also has an old "nmc" currency
		// in the orderbook, we otta ignore the nmc currency
		if book, ok := e.books[orders[i].currency]; ok {
			book.insert(orders[i])
		}
	}
}

func fillBookWithFakeOrders(book *orderbook, currency string) {

	var last int64 = 100
	for i := 0; i < 10; i++ {
		order := createOrder("joshua", rand.Int63n(200)+1, last, SELL, currency)
		results := book.match(order)
		err := storeOrder(order, results)
		if err != nil {
			panic(err)
		}
		last += rand.Int63n(5) + 1
	}

	last = 99
	for i := 0; i < 10; i++ {
		order := createOrder("jeffery", rand.Int63n(200)+1, last, BUY, currency)
		results := book.match(order)
		err := storeOrder(order, results)
		if err != nil {
			panic(err)
		}
		last -= rand.Int63n(5) + 1
	}
}
