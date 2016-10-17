// What is an execution one might ask?
// This exchange operates in a limbo between memory and a SQL database,
// essentially keeping the whole orderbook cached in memory and persisted
// on disk via database. An execution is essentially an update to the state
// of the whole database and also an update to the orderbook. Basically an
// execution translates to SQL insert row or update balance
package main

import (
	"time"
)

const (
	DEPTH = 10
)

// the following constants represent how an order has been
// executed, and what should happen in the database
const (
	PARTIAL = iota // a partial fill means we should subtract from an order in DB
	FULL           // a full fill means we should remove an order from DB
	OPEN           // an open order means we've created a whole new order, should get inserted
	CANCEL         // user decided to cancel own order
)

// execution represents an order that was matched and executed.
// It uses PriceSum because it's nicer to just add together order
// prices instead of figuring out the average price.
type execution struct {
	ID         int64
	Name       string
	Filler     string
	Amount     int64
	Price      int64
	Order_type int
	Status     int
	Currency   string
	Timestamp  time.Time
}

// executions is a history of order executions for an orderbook
// uses a circular array for performance additions
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
