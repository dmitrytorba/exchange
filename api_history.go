package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func historyHandler(w http.ResponseWriter, r *http.Request) error {
	username := "somerandomguy" // placeholder, waiting for an authentication system to supply this bit of information

	sel, err := db.Prepare("SELECT id, amount, price, order_type, filler, username, currency, timestamp FROM executions WHERE username = $1 OR filler = $1")
	if err != nil {
		return err
	}

	rows, err := sel.Query(username)
	if err != nil {
		return err
	}
	execs := make([]*execution, 0, 32)

	defer rows.Close()
	for rows.Next() {
		execution := &execution{}
		var order_type string
		var timestamp string
		if err := rows.Scan(&execution.ID, &execution.Amount, &execution.Price, &order_type, &execution.Filler, &execution.Name, &execution.Currency, &timestamp); err != nil {
			return err
		}

		if order_type == "buy" {
			execution.Order_type = BUY
		} else {
			execution.Order_type = SELL
		}

		execution.Timestamp, err = time.Parse("2006-01-02T15:04:05.999999Z", timestamp)
		if err != nil {
			return err
		}

		execs = append(execs, execution)
	}

	return json.NewEncoder(w).Encode(execs)
}
