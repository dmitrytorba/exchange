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

	return e
}

func fillBookWithFakeOrders(book *orderbook) {

	var last int64 = 100
	for i := 0; i < 30; i++ {
		order := createOrder("joshua", rand.Int63n(200)+1, last, SELL)
		book.history.addExecution(&execution{
			Name:       order.Name,
			Amount:     order.Amount,
			PriceSum:   order.Price * order.Amount,
			Order_type: SELL,
			Status:     OPEN,
		})
		book.insert(order)
		last += rand.Int63n(5) + 1
	}

	last = 99
	for i := 0; i < 30; i++ {
		order := createOrder("jeffery", rand.Int63n(200)+1, last, BUY)
		book.history.addExecution(&execution{
			Name:       order.Name,
			Amount:     order.Amount,
			PriceSum:   order.Price * order.Amount,
			Order_type: BUY,
			Status:     OPEN,
		})
		book.insert(order)
		last -= rand.Int63n(5) + 1
	}
}
