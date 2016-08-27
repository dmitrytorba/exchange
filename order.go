package main

import ()

type currency string

const ( // enumerated types
	BTC currency = "BTC"
	USD          = "USD"
	LTC          = "LTC"
)

const (
	SELL int = iota
	BUY      = iota
)

type order struct {
	OrderType int `json:"type"`
	Amount    int
	User      string
	Price     int
}
