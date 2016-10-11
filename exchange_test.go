package main

import (
	"testing"
)

// thoe following bechmark will not use the SQL db
func BenchmarkOrderInsert(b *testing.B) {
	startDb() // only used to satisfy the start of exchange
	startExchange()
	b.StartTimer()
	currency := exch.currencies[0]
	for n := 0; n < b.N; n++ {
		exch.books[currency].match(createOrder("jay", 1, int64(n), SELL, currency))
	}
}
