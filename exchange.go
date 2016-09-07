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

	e.books["ltc"] = createOrderbook()
	e.books["eth"] = createOrderbook()
	fillBookWithFakeOrders(e.books["ltc"])

	return e
}

func fillBookWithFakeOrders(book *orderbook) {

	var last int64 = 100
	for i := 0; i < 30; i++ {
		book.insert(createOrder("joshua", 100, last, SELL))
		last += rand.Int63n(5) + 1
	}

	last = 99
	for i := 0; i < 30; i++ {
		book.insert(createOrder("jeffery", 100, last, BUY))
		last -= rand.Int63n(5) + 1
	}
}
