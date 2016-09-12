package main

import ()

const (
	DEPTH = 100
)

// execution represents an order that was matched and executed
type execution struct {
	Name       string
	Amount     int64
	Price      int64
	Order_type int
	Status     string
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
