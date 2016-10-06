package main

import (
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ob := exch.books["ltc"]
	buys := ob.array(BUY)
	sells := ob.array(SELL)

	executeTemplate(w, "home", 200, map[string]interface{}{
		"Currencies": exch.currencies,
		"Sells":      sells,
		"Buys":       buys,
		"Executions": ob.history.array(),
	})
}
