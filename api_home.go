package main

import (
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := checkMe(r)
	if err != nil {
		return err
	}

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

	if user != nil {
		data["Username"] = user.username
		data["Authed"] = true
	}

	return executeTemplate(w, "home", 200, data)
}
