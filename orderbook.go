package main

import (
	"container/list"
	"sync"
	"time"
)

// the following constants represent the type of an order
const (
	BUY = iota
	SELL
)

type orderbook struct {
	mtx     sync.RWMutex
	buys    *list.List
	sells   *list.List
	history *executions
	ticker  chan *execution
}

// my idea for consistency:
// keep dual copies of orderbook on sql and on memory.
// keep orders unordered on psql
// execute order using in memory orderbook.
// when its time to update balances, also update/delete orders in
// one good transaction.

// order represents an open order
type order struct {
	ID         int64
	Name       string
	Amount     int64
	Price      int64
	order_type int
	currency   string
	timestamp  time.Time
}

func createOrder(name string, amount, price int64, order_type int, currency string) *order {
	return &order{
		Name:       name,
		Amount:     amount,
		Price:      price,
		order_type: order_type,
		currency:   currency,
		timestamp:  time.Now(),
	}
}

func createOrderbook() *orderbook {
	ob := &orderbook{
		sells:   list.New(),
		buys:    list.New(),
		history: createExecutions(),
		ticker:  make(chan *execution),
	}
	return ob
}

// array takes the linked list and spits it out as an array
func (o *orderbook) array(order_type int) []*order {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	var array []*order
	if order_type == BUY {
		array = make([]*order, o.buys.Len()) // Note .Len() is not a typo that came accidentally from JS
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
// the orderbook. It is thread-safe.
func (o *orderbook) match(initial *order) []*execution {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	// collect all the orders we executed
	execs := make([]*execution, 0, 10)

	// get the opposite list to match against
	// and also a pretty string to show in the executions table
	list := o.sells
	if initial.order_type == SELL {
		list = o.buys
	}

	var matched *order
	e := list.Front()
	for e != nil {
		matched = e.Value.(*order)
		e = e.Next()

		// if matching order is a buy and price is below the buy order, FILL!
		if (initial.order_type == BUY && matched.Price <= initial.Price) || (initial.order_type == SELL && matched.Price >= initial.Price) {
			if initial.Amount >= matched.Amount { // matching order is overfilled, we must remove it
				exec := &execution{
					ID:         matched.ID,
					Name:       matched.Name,
					Filler:     initial.Name,
					Amount:     matched.Amount,
					Price:      matched.Price,
					Order_type: matched.order_type,
					Status:     FULL,
					Currency:   matched.currency,
					Timestamp:  time.Now(),
				}
				// add our execution to the ticker and history
				o.ticker <- exec
				execs = append(execs, exec)
				o.history.addExecution(exec)

				// remove from the list, being careful to set the next iteration
				list.Remove(e.Prev())

				// remove from the initial order since it has not been filled yet
				initial.Amount -= matched.Amount

			} else { // matching order fills initial order fully

				exec := &execution{
					ID:         matched.ID,
					Name:       matched.Name,
					Filler:     initial.Name,
					Amount:     initial.Amount,
					Price:      matched.Price,
					Order_type: matched.order_type,
					Status:     PARTIAL,
					Currency:   matched.currency,
					Timestamp:  time.Now(),
				}

				// add our execution to the ticker and history
				o.ticker <- exec
				execs = append(execs, exec)
				o.history.addExecution(exec)

				// decrease the matched order in order to fill initial order
				matched.Amount -= initial.Amount
				initial.Amount = 0
			}
		} else { // if no matching order can be executed, shelve the order to be executed later
			break
		}

		// It's a good idea to stop filling orders if your matching order has been 100% filled
		if initial.Amount == 0 {
			break
		}
	}

	// if the order is not fully filled its a hot idea to insert it into the orderbook
	// to later get filled
	if initial.Amount > 0 {
		o.insert(initial)
	}

	return execs
}

// insert will insert an order into an orderbook at the correct
// position. It could be improved to avoid looping
// through the entire list to place the last element.
// This is not thread-safe.
func (o *orderbook) insert(addOrder *order) {
	// get the list the order belongs to
	list := o.sells
	if addOrder.order_type == BUY {
		list = o.buys
	}

	var matched *order
	for e := list.Front(); e != nil; e = e.Next() {
		matched = e.Value.(*order)

		// lower priced orders first for sells
		if addOrder.order_type == SELL && addOrder.Price < matched.Price {
			list.InsertBefore(addOrder, e)
			return
		}

		// higher priced orders first for buys
		if addOrder.order_type == BUY && addOrder.Price > matched.Price {
			list.InsertBefore(addOrder, e)
			return
		}
	}

	list.PushBack(addOrder)
}

func (o *orderbook) getLeadBuyPrice() int64 {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	if o.buys.Len() == 0 {
		return 0
	}
	
	var leader *order = o.buys.Front().Value.(*order)
	return leader.Price
}

func (o *orderbook) getLeadSellPrice() int64 {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	if o.sells.Len() == 0 {
		return 0
	}
	
	var leader *order = o.sells.Front().Value.(*order)
	return leader.Price
}

func (o *orderbook) readTicker() {
	for {
		<-o.ticker
	}
}
