package main

// insertIndex figures out where in the list the order fits into
func insertIndex(orders []*order, order *order) int {
	if len(orders) == 0 {
		return 0
	}

	for i := 0; i < len(orders); i++ {
		compare := orders[i]

		if order.OrderType == BUY && compare.Price > order.Price {
			return i
		}

		if order.OrderType == SELL && compare.Price < order.Price {
			return i
		}
	}

	return len(orders)
}
