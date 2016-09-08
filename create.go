package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func orderHandler(w http.ResponseWriter, r *http.Request) {
	amount := r.FormValue("amount")
	price := r.FormValue("price")
	action := r.FormValue("action")

	pricei, err := strconv.ParseInt(price, 10, 64)
	if err != nil {
		executeTemplate(w, "error", 200, map[string]interface{}{})
		return
	}

	amounti, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		executeTemplate(w, "error", 200, map[string]interface{}{})
		return
	}

	var execs []*execution
	if action == "buy" {
		execs = exch.books["ltc"].match(createOrder("somerandomguy", amounti, pricei, BUY))
	} else if action == "sell" {
		execs = exch.books["ltc"].match(createOrder("somerandomguy", amounti, pricei, SELL))
	} else {
		executeTemplate(w, "error", 200, map[string]interface{}{})
		return
	}

	fmt.Println("matched", len(execs), "orders")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
