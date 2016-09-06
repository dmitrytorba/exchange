package main

import ()

type exchange struct {
}

func createExchange() (*exchange, error) {
	return &exchange{}, nil
}
