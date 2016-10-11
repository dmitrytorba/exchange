// this file will ideally hook up to websockets and tick out live price
// updates for our users
package main

import ()

const (
	MESSAGE_CHUNK = 100
)

type ticker struct {
	ticks chan *execution
}

func (t *ticker) addExecution(exec *execution) {
	t.ticks <- exec
}

func StartTicker() *ticker {
	t := &ticker{
		ticks: make(chan *execution),
	}

	return t
}
