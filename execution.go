package main

import ()

const (
	DEPTH = 100
)

// the following constants represent how an order has been
// executed, and what should happen in the database
const (
	PARTIAL = iota // a partial fill means we should subtract from an order in DB
	FULL           // a full fill means we should remove an order from DB
	OPEN           // an open order means we've created a whole new order, should get inserted
)

// execution represents an order that was matched and executed.
// I used PriceSum because it's nicer to just add together order
// prices instead of figuring out the average price.
type execution struct {
	Name       string
	Amount     int64
	PriceSum   int64
	Order_type int
	Status     int
}

// executions is a history of order executions for an orderbook
type executions struct {
	history []*execution
	start   int
	length  int
}

func createExecutions() *executions {
	return &executions{
		start:   0,
		length:  0,
		history: make([]*execution, DEPTH),
	}
}

func (e *executions) addExecution(exec *execution) {
	e.history[(e.start+e.length)%DEPTH] = exec

	if e.length == DEPTH {
		e.start = (e.start + 1) % DEPTH
	} else {
		e.length++
	}
}

// 0

func (e *executions) array() []*execution {
	array := make([]*execution, e.length)

	for i := 0; i < e.length; i++ {
		array[i] = e.history[(e.start+e.length-i-1)%DEPTH]
	}

	return array
}
