// this was my first naive try at making an automated trader
// the strategy was to simply match the existing highest/lowest
// sell/buy orders on the order book plus some delta 'buffer'
// no orders are actually placed; the code simulates a fill if
// theres a trade executed with a matching price 

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

var btc float64
var usd float64

const tick float64 = 0.1
const targetDelta float64 = 1

func round(f float64) float64 {
	return math.Floor(f*10 + 0.5)/10
}

func StartTrader() {
	myBid = 0
	myAsk = 0
	bidDelta = targetDelta
	askDelta = targetDelta
	btc = 1
	usd = 1100
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
	if myBid == 0  || myBid != targetBid {
		myBid = targetBid
		log.Printf("bid: %f (delta=%f)", myBid, bidDelta)
	}

	targetAsk := round(ask.Price + (askDelta*tick))
	if myAsk == 0  || myAsk != targetAsk {
		myAsk = targetAsk
		log.Printf("ask: %f (delta=%f)", myAsk, askDelta)
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
	//log.Printf("bitfinex trade, price: %s, vol %s", price, volume)

	if volume > 0 {
		if price >= myAsk {
			btc -= 0.1
			usd += 0.1*price
			fee := 0.0001*price
			usd -= fee
			log.Printf("my ask filled @ %f (fee=%f) [ balance now btc=%f  usd=%f ]", myAsk, fee, btc, usd)
			if bidDelta > targetDelta {
				bidDelta--
			} else {
				askDelta++
			}
			myAsk=0
			updateOrders()
		}
	} else {
		if price <= myBid {
			btc += 0.1
			usd -= 0.1*price
			fee := 0.0001*price
			usd -= fee
			log.Printf("my bid filled @ %f (fee=%f) [ balance now btc=%f  usd=%f ]", myBid, fee, btc, usd)
			if askDelta > targetDelta {
				askDelta--
			} else {
				bidDelta++
			}
			myBid = 0
			updateOrders()
		}
	}
}
