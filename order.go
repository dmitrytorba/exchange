package main

import (
	"net/http"
	"strconv"
)

func orderHandler(w http.ResponseWriter, r *http.Request) {
	amount := r.FormValue("amount")
	price := r.FormValue("price")
	action := r.FormValue("action")

	pricei, err := strconv.ParseInt(price, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	amounti, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if action == "buy" {
		exch.books["ltc"].match(createOrder("somerandomguy", amounti, pricei, BUY))
	} else if action == "sell" {
		exch.books["ltc"].match(createOrder("somerandomguy", amounti, pricei, SELL))
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
