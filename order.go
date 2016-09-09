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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	amounti, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var execs []*execution
	if action == "buy" {
		execs = exch.books["ltc"].match(createOrder("somerandomguy", amounti, pricei, BUY))
	} else if action == "sell" {
		execs = exch.books["ltc"].match(createOrder("somerandomguy", amounti, pricei, SELL))
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// adding the user's order to our execution table
	execs = append(execs, &execution{
		Name:   "(you)",
		Price:  pricei,
		Amount: amounti,
		Type:   fmt.Sprintf("initial %v", action),
		Status: "PROCESSED",
	})
	exch.recent = append(execs, exch.recent...)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
