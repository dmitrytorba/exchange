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

	// check the db to see if its empty, if it is fill it with
	// fake orders fo dev purposes
	var count int64
	err := db.QueryRow("SELECT count(*) FROM orders").Scan(&count)
	if err != nil {
		panic(err)
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
	// also fill with fake orders if needed
	for i := 0; i < len(e.currencies); i++ {
		currency := e.currencies[i]
		e.books[currency] = createOrderbook()
		if count == 0 {
			fillBookWithFakeOrders(e.books[currency], currency)
		}
	}

	// fill the orderbook with persisted orders
	if count > 0 {
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

	return e
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
