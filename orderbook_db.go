package main

import (
	"fmt"
)

const (
	DEFAULT_CURRENCY = "btc"
)

// getAllOrders will query all rows from the db,
// it is intended to use as a recovery process
// for getting the orderbook synchronized with
// the in database
func getAllOrders() ([]*order, error) {
	orders := make([]*order, 0, 100)
	rows, err := db.Query("SELECT id, amount, price, order_type, username, currency FROM orders")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		order := &order{}
		var order_type string
		if err := rows.Scan(&order.ID, &order.Amount, &order.Price, &order_type, &order.Name, &order.currency); err != nil {
			return nil, err
		}
		if order_type == "buy" {
			order.order_type = BUY
		} else {
			order.order_type = SELL
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// storeOrder will take an outstanding order or nil, along with
// it's exectutions (or nil) and persist them in the database
func storeOrder(order *order, execs []*execution) error {

	// start the transaction and prep the statements
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	remove, err := tx.Prepare(`DELETE FROM orders WHERE id=$1`)
	if err != nil {
		return err
	}
	update, err := tx.Prepare(`UPDATE orders SET amount = amount - $1 WHERE id=$2`)
	if err != nil {
		return err
	}

	defaultbal, err := tx.Prepare(fmt.Sprintf(`UPDATE users SET %v=%v + $1 WHERE username=$2`, DEFAULT_CURRENCY, DEFAULT_CURRENCY))
	if err != nil {
		return err
	}

	var othercurrency string
	if len(execs) > 0 {
		othercurrency = execs[0].Currency
	} else if order != nil {
		othercurrency = order.currency
	}

	otherbal, err := tx.Prepare(fmt.Sprintf(`UPDATE users SET %v=%v + $1 WHERE username=$2`, othercurrency, othercurrency))
	if err != nil {
		return err
	}

	// loop through the executions and do corresponding updates or removes
	for i := 0; i < len(execs); i++ {
		exec := execs[i]
		var err error
		switch exec.Status {
		case PARTIAL:
			_, err = update.Exec(exec.Amount, exec.ID)
		case FULL:
			_, err = remove.Exec(exec.ID)
		}

		if err != nil {
			return err
		}

		// handle per execution balance logic here
		// NOTE: the only person who should get his balance decreased is always you
		// the original owner of the order has already had his balance decreased for either the
		// target currency or the default currency
		if exec.Order_type == BUY {
			// the creator of this order being executed is looking to buy the currency and has already
			// payed with the default currency, meanwhile you are looking to sell your other currency and
			// get the default currency
			otherbal.Exec(exec.Amount, exec.Name)
			otherbal.Exec(-exec.Amount, exec.Filler)

			defaultbal.Exec(exec.Amount*exec.Price, exec.Filler)
		} else {
			// the creator of this order being executed is looking to sell this currency (already payed with it) and receive payment for
			// it in the default currency, while you are looking to buy his currency and pay with the default currency
			otherbal.Exec(exec.Amount, exec.Filler)

			defaultbal.Exec(exec.Amount*exec.Price, exec.Name)
			defaultbal.Exec(-exec.Amount*exec.Price, exec.Filler)
		}
	}

	// insert the order into the orderbook
	if order != nil && order.Amount > 0 {
		insert, err := tx.Prepare(`INSERT INTO orders (amount, price, order_type, username, currency) VALUES ($1, $2, $3, $4, $5) RETURNING id`)
		if err != nil {
			return err
		}

		typestring := "sell"
		if order.order_type == BUY {
			typestring = "buy"
		}

		err = insert.QueryRow(order.Amount, order.Price, typestring, order.Name, order.currency).Scan(&order.ID)
		if err != nil {
			return err
		}

		// handle final order balance movement here
		if order.order_type == BUY {
			// You are looking to buy the currency and pay with the default currency, we're gonna make you
			// preemptively pay in the default currency so you can't just put down an order and then withdraw
			// your whole balance
			defaultbal.Exec(-order.Amount*order.Price, order.Name)
		} else {
			// your looking to receive the default currency for these coins that you have, so were gonna
			// do the opposite of what I said above but we're still gonna debit you in the same way
			otherbal.Exec(order.Amount, order.Name)
		}
	}

	return tx.Commit()

}
