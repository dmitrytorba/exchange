package main

import (
	"math/rand"
)

type exchange struct {
	books map[string]*orderbook
}

func createExchange() *exchange {
	e := &exchange{
		books: make(map[string]*orderbook),
	}

	// check the db to see if its empty, if it is fill it with
	// fake orders fo dev purposes
	var count int64
	err := db.QueryRow("SELECT count(*) FROM orders").Scan(&count)
	if err != nil {
		panic(err)
	}

	// cycle through the config currencies, creating an orderbook for each
	// also fill with fake orders if needed
	for i := 0; i < len(currencies); i++ {
		currency := currencies[i]
		e.books[currencies[i]] = createOrderbook()

		if count == 0 {
			fillBookWithFakeOrders(e.books[currency], currency)
		}
	}

	// fill the orderbook with persisted orders
	orders, err := getAllOrders()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(orders); i++ {
		e.books[orders[i].currency].insert(orders[i])
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
