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

	var executions []*execution
	var order *order
	if action == "buy" { // note that the "ltc" designation is temporary
		order = createOrder("somerandomguy", amounti, pricei, BUY, "ltc")
		executions = exch.books["ltc"].match(order)
	} else if action == "sell" {
		order = createOrder("somerandomguy", amounti, pricei, SELL, "ltc")
		executions = exch.books["ltc"].match(order)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = storeOrder(order, executions)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
