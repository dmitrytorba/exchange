package main

import ()

func getAllOrders() ([]*order, error) {
	orders := make([]*order, 0, 100)
	rows, err := db.Query("SELECT id, amount, price, order_type, username FROM orders")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		order := &order{}
		var order_type string
		if err := rows.Scan(&order.ID, &order.Amount, &order.Price, &order_type, &order.Name); err != nil {
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

func storeOrder(order *order, execs []*execution) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	remove, err := tx.Prepare(`DELETE FROM orders WHERE id=$1`)
	if err != nil {
		return err
	}

	insert, err := tx.Prepare(`INSERT INTO orders (amount, price, order_type, username) VALUES ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		return err
	}

	update, err := tx.Prepare(`UPDATE orders SET amount = amount - $1 WHERE id=$2`)
	if err != nil {
		return err
	}

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
	}

	if order != nil {
		typestring := "sell"
		if order.order_type == BUY {
			typestring = "buy"
		}

		err = insert.QueryRow(order.Amount, order.Price, typestring, order.Name).Scan(&order.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()

}
