// this file will ideally hook up to websockets and tick out live price
// updates for our users
package main

import (
	"fmt"
)

const (
	MESSAGE_THRESHOLD = 100
)

type ticker struct {
	ticks chan []*execution
}

func (t *ticker) readTicks() {
	for i := range t.ticks {
		fmt.Println(i)
	}
}

func StartTicker() *ticker {
	t := &ticker{
		ticks: make(chan []*execution),
	}

	go t.readTicks()

	return t
}
