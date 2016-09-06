package main

import (
	"container/list"
	"fmt"
	"time"
)

const (
	BUY = iota
	SELL
)

type orderbook struct {
	buys  *list.List
	sells *list.List
}

// order represents an open order
type order struct {
	name       string
	amount     int64
	price      int64
	order_type int
	open       bool
	timestamp  time.Time
}

type execution struct {
	name   string
	amount int64
	price  int64
}

func createOrder(name string, amount, price int64, order_type int) *order {
	return &order{
		name:       name,
		amount:     amount,
		price:      price,
		order_type: order_type,
		timestamp:  time.Now(),
		open:       true,
	}
}

func createOrderbook() *orderbook {
	ob := &orderbook{}
	ob.sells = list.New()
	ob.buys = list.New()
	return ob
}

func (o *orderbook) match(matchOrder *order) []*execution {
	// collect all the orders we executed
	execs := make([]*execution, 0, 10)

	// get the opposite list to match against
	list := o.sells
	if matchOrder.order_type == SELL {
		list = o.buys
	}

	var iter *order
	for e := list.Front(); e != nil; e = e.Next() {
		iter = e.Value.(*order)

		// if matching order is a buy and price is below the buy order, FILL!
		if (matchOrder.order_type == BUY && iter.price <= matchOrder.price) || (matchOrder.order_type == SELL && iter.price >= matchOrder.price) {
			if matchOrder.amount >= iter.amount {
				// remove the order, it has been filled
				execs = append(execs, &execution{
					name:   iter.name,
					amount: iter.amount,
					price:  iter.price,
				})
				list.Remove(e)
				matchOrder.amount -= iter.amount
			} else if matchOrder.amount < iter.amount {
				execs = append(execs, &execution{
					name:   iter.name,
					amount: matchOrder.amount,
					price:  iter.price,
				})
				iter.amount -= matchOrder.amount
				matchOrder.amount = 0
			}
		} else { // if no matching order can be executed, shelve the order to be executed later
			break
		}

		// It's a good idea to stop filling orders if your matching order has been 100% filled
		if matchOrder.amount == 0 {
			break
		}

	}

	if matchOrder.amount > 0 {
		o.insert(matchOrder)
	}

	return execs
}

func (o *orderbook) insert(addOrder *order) {
	// get the list the order belongs to
	list := o.sells
	if addOrder.order_type == BUY {
		list = o.buys
	}

	var iter *order
	for e := list.Front(); e != nil; e = e.Next() {
		iter = e.Value.(*order)

		// lower priced orders first for sells
		if addOrder.order_type == SELL && addOrder.price < iter.price {
			list.InsertBefore(addOrder, e)
			return
		}

		// higher priced orders first for buys
		if addOrder.order_type == BUY && addOrder.price > iter.price {
			list.InsertBefore(addOrder, e)
			return
		}
	}

	list.PushBack(addOrder)
}

func (o *orderbook) print() {
	fmt.Println("SELLS")
	fmt.Println("======================")
	for e := o.sells.Front(); e != nil; e = e.Next() {
		order := e.Value.(*order)
		fmt.Printf("$%v/%v-%v\n", order.price, order.amount, order.name)
	}

	fmt.Println("BUYS")
	fmt.Println("======================")
	for e := o.buys.Front(); e != nil; e = e.Next() {
		order := e.Value.(*order)
		fmt.Printf("$%v/%v-%v\n", order.price, order.amount, order.name)
	}
}
