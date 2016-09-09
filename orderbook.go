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
	Name       string
	Amount     int64
	Price      int64
	order_type int
	timestamp  time.Time
}

// execution represents an order that was matched and executed
type execution struct {
	Name   string
	Amount int64
	Price  int64
	Type   string
	Status string
}

func createOrder(name string, amount, price int64, order_type int) *order {
	return &order{
		Name:       name,
		Amount:     amount,
		Price:      price,
		order_type: order_type,
		timestamp:  time.Now(),
	}
}

func createOrderbook() *orderbook {
	ob := &orderbook{}
	ob.sells = list.New()
	ob.buys = list.New()
	return ob
}

// array takes the linked list and spits it out as an array
func (o *orderbook) array(order_type int) []*order {
	var array []*order
	if order_type == BUY {
		array = make([]*order, o.buys.Len())
	} else {
		array = make([]*order, o.sells.Len())
	}

	if len(array) == 0 {
		return array
	}

	e := o.sells.Front()
	if order_type == BUY {
		e = o.buys.Front()
	}
	for i := 0; i < len(array); i++ {
		array[i] = e.Value.(*order)
		e = e.Next()
	}
	return array
}

// match will match orders, deleting/modifying orders that get executed.
// once matching is finished, the order, if unfilled, will be inserted into
// the orderbook.
func (o *orderbook) match(matchOrder *order) []*execution {
	// collect all the orders we executed
	execs := make([]*execution, 0, 10)

	// get the opposite list to match against
	// and also a pretty string to show in the executions table
	list := o.sells
	pretty_type := "sell"
	if matchOrder.order_type == SELL {
		list = o.buys
		pretty_type = "buy"
	}
	pretty_type = fmt.Sprintf("matched %v", pretty_type)

	var iter *order
	e := list.Front()
	for e != nil {
		iter = e.Value.(*order)

		// if matching order is a buy and price is below the buy order, FILL!
		if (matchOrder.order_type == BUY && iter.Price <= matchOrder.Price) || (matchOrder.order_type == SELL && iter.Price >= matchOrder.Price) {

			if matchOrder.Amount >= iter.Amount { // matching order is overfilled, we must remove it
				// remove the order, it has been filled
				execs = append(execs, &execution{
					Name:   iter.Name,
					Amount: iter.Amount,
					Price:  iter.Price,
					Type:   pretty_type,
					Status: "FULL EXECUTION",
				})
				e = e.Next()
				list.Remove(e.Prev())
				matchOrder.Amount -= iter.Amount
			} else { // matching order fills initial order fully
				execs = append(execs, &execution{
					Name:   iter.Name,
					Amount: matchOrder.Amount,
					Price:  iter.Price,
					Type:   pretty_type,
					Status: "PARTIAL EXECUTION",
				})
				iter.Amount -= matchOrder.Amount
				matchOrder.Amount = 0
				e = e.Next()
			}
		} else { // if no matching order can be executed, shelve the order to be executed later
			break
		}

		// It's a good idea to stop filling orders if your matching order has been 100% filled
		if matchOrder.Amount == 0 {
			break
		}

	}

	if matchOrder.Amount > 0 {
		o.insert(matchOrder)
	}

	return execs
}

// insert will insert an order into an orderbook at the correct
// position. It could be improved to avoid looping
// through the entire list to place the last element.
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
		if addOrder.order_type == SELL && addOrder.Price < iter.Price {
			list.InsertBefore(addOrder, e)
			return
		}

		// higher priced orders first for buys
		if addOrder.order_type == BUY && addOrder.Price > iter.Price {
			list.InsertBefore(addOrder, e)
			return
		}
	}

	list.PushBack(addOrder)
}
