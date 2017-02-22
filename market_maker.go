package main

import (
	"log"
	"encoding/json"
	"math"
)

var myBid float64
var myAsk float64

var bidDelta float64
var askDelta float64

const tick float64 = 0.1

func round(f float64) float64 {
	return math.Floor(f*10 + 0.5)/10
}

func StartTrader() {
	myBid = 0
	myAsk = 0
	bidDelta = 0
	askDelta = 0
	pubsub, err := rd.Subscribe("bitfinex-btcusd")
	if err != nil {
    panic(err)
	}
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}
		if msg.Payload == "book" {
			updateOrders()
		}
		if msg.Payload == "trade" {
			onTrade()
		}
	}
}

func updateOrders() {
	key := "btcusd-bitfinex-asks"
	vals, err := rd.ZRange(key, 0, 0).Result()
	
	if err != nil || len(vals)  != 1 {
		log.Fatal("redis err: ", err)
	}
	entryStr := vals[0]
	var ask BitfinexBookEntry 
	json.Unmarshal([]byte(entryStr), &ask)
	
	key = "btcusd-bitfinex-bids"
	vals, err = rd.ZRange(key, -1, -1).Result()
	
	if err != nil || len(vals)  != 1 {
		log.Fatal("redis err: ", err)
	}
	entryStr = vals[0]
	var bid BitfinexBookEntry 
	json.Unmarshal([]byte(entryStr), &bid)

	// log.Printf("bitfinex book, bid: %s, ask %s", bid.Price, ask.Price)

	targetBid := round(bid.Price - (bidDelta*tick))
	if myBid == 0  || myBid < targetBid {
		myBid = targetBid
		log.Printf("orders: %s ... %s", myBid, myAsk)
	}

	targetAsk := round(ask.Price + (askDelta*tick))
	if myAsk == 0  || myAsk > targetAsk {
		myAsk = targetAsk
		log.Printf("orders: %s ... %s", myBid, myAsk)
	}
}

func onTrade() {
	var price float64
	var volume float64
	queryStr := "select price, volume from bitfinex_trades_btcusd order by time_stamp desc limit 1"
	err := db.QueryRow(queryStr).Scan(&price, &volume)
	if err != nil {
		log.Fatal("trade select err", err)
	}
	log.Printf("bitfinex trade, price: %s, vol %s", price, volume)

	if volume > 0 {
		if price == myAsk {
			log.Printf("my ask filled! %s", price)
			askDelta++
			myAsk=0
			updateOrders()
		}
	} else {
		if price == myBid {
			log.Printf("my bid filled! %s", price)
			bidDelta++
			myBid = 0
			updateOrders()
		}
	}
}
