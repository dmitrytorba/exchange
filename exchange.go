package main

import ()

type exchange struct {
	books map[string]*orderbook
}

func createExchange() *exchange {
	e := &exchange{
		books: make(map[string]*orderbook),
	}

	e.books["ltc"] = createOrderbook()
	e.books["eth"] = createOrderbook()

	return e
}
