package main

import (
	"time"
)

type exchange struct {
}

func createExchange() (*exchange, error) {
	return &exchange{}, nil
}

type order struct {
	id         int
	user_id    int
	price      int64
	amount     int64
	currency   string
	order_type string
	timestamp  *time.Time
}

func (e *exchange) createOrder(user_id int, price int64, amount int64, currency string, order_type string) (*order, error) {
	order := &order{
		user_id:    user_id,
		price:      price,
		amount:     amount,
		currency:   currency,
		order_type: order_type,
	}

	// Begin this sql transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	insert, err := tx.Prepare("INSERT INTO orders(order_type, currency, amount, price, user_id) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		panic(err)
	}

	_, err = insert.Exec(order_type, currency, amount, price, user_id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return order, nil
}
