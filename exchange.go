package main

import (
	"fmt"
	"math/rand"
)

type exchange struct {
	books  map[string]*orderbook
	recent []*execution
}

func createExchange() *exchange {
	e := &exchange{
		books:  make(map[string]*orderbook),
		recent: make([]*execution, 0, 100),
	}

	e.books["ltc"] = createOrderbook()
	e.books["eth"] = createOrderbook()
	fillBookWithFakeOrders(e.books["ltc"])

	return e
}

func fillBookWithFakeOrders(book *orderbook) {

	execs := make([]*execution, 30)
	var last int64 = 100
	for i := 0; i < 30; i++ {
		order := createOrder("joshua", rand.Int63n(200)+1, last, SELL)
		results := book.match(order)
		execs[i] = results[0]
		last += rand.Int63n(5) + 1
	}
	fmt.Println(processExecutions(execs))

	last = 99
	for i := 0; i < 30; i++ {
		order := createOrder("jeffery", rand.Int63n(200)+1, last, BUY)
		results := book.match(order)
		execs[i] = results[0]
		last -= rand.Int63n(5) + 1
	}
	fmt.Println(processExecutions(execs))
}
