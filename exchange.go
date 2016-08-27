package main

import ()

type exchange struct {
}

func createExchange() (*exchange, error) {
	e := &exchange{}

	return e, nil
}
