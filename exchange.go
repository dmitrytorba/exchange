package main

import (
	"fmt"
	"time"
)

type currency string

const ( // enumerated types
	BTC currency = "BTC"
	USD          = "USD"
	LTC          = "LTC"
)

type exchange struct {
	reserves map[currency]int
}

func createExchange() *exchange {
	exch := &exchange{}

	exch.reserves = make(map[currency]int)
	exch.reserves[USD] = 100
	exch.reserves[BTC] = 100

	return exch
}

func (e *exchange) execute(amount int, to, from currency) {
	// the pricing algorithm, 1:1 for now
	conversion := 1

	e.reserves[from] += amount
	e.reserves[to] -= amount * conversion
}

type order struct {
	amount int
	from   currency
	to     currency
}

func (e *exchange) printout() {
	fmt.Println("ex.ch.an.ge ------- ", time.Now().Format("Mon Jan 2 15:04:05"))

	fmt.Print("\n")

	fmt.Println("  currency   |     amount   ")
	fmt.Println("----------------------------")
	for k, v := range e.reserves {
		fmt.Printf("    %v             %v\n", k, v)
	}
}
