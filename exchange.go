package main

import (
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

	orders, err := getAllOrders()
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(orders); i++ {
		e.books["ltc"].insert(orders[i])
	}

	return e
}

func fillBookWithFakeOrders(book *orderbook) {

	var count int64
	err := db.QueryRow("SELECT count(*) FROM orders").Scan(&count)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		var last int64 = 100
		for i := 0; i < 30; i++ {
			order := createOrder("joshua", rand.Int63n(200)+1, last, SELL)
			results := book.match(order)
			err = storeOrder(order, results)
			if err != nil {
				panic(err)
			}
			last += rand.Int63n(5) + 1
		}

		last = 99
		for i := 0; i < 30; i++ {
			order := createOrder("jeffery", rand.Int63n(200)+1, last, BUY)
			results := book.match(order)
			err = storeOrder(order, results)
			if err != nil {
				panic(err)
			}
			last -= rand.Int63n(5) + 1
		}
	}
}
