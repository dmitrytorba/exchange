package main

import (
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ob := exch.books["ltc"]
	buys := ob.array(BUY)
	sells := ob.array(SELL)

	data := map[string]interface{}{
		"Currencies": exch.currencies,
		"Sells":      sells,
		"Buys":       buys,
		"Executions": ob.history.array(),
		"LeadBuy":    ob.getLeadBuyPrice(),
		"LeadSell":   ob.getLeadSellPrice(),
	}

	usr := getUserFromCookie(r)
	if usr != nil {
		data["Username"] = usr.username
	}

	executeTemplate(w, "home", 200, data)
}
