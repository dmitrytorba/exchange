package main

import (
	"log"
	"encoding/json"
	"strconv"
	"time"
)

const gdaxWS ="wss://ws-feed.gdax.com"

func connectGdax() {
	monitorWebsocket(
		gdaxWS,
		`{"type":"subscribe","product_ids":["BTC-USD"]}`,
		onGdaxEvent)
}

type GdaxMsg struct {
	Type string
	Side string
	OrderId string `json:"order_id"`
	Reason string
	ProductId string
	Size string
	Price string
	RemainingSize string `json:"remaining_size"`
	Sequence int64
	Time string
	OrderType string
	Funds string
	MakerOrderId string
	TakerOrderId string
	NewSize string `json:"new_size"`
	OldSize string
	NewFunds string
	OldFunds string
	LastTradeId string
	Message string
}

func onGdaxEvent(messageStr string) {
	var msg GdaxMsg
	json.Unmarshal([]byte(messageStr), &msg)
	switch msg.Type {
	case "error":
		log.Printf("GDAX error: %s", msg.Message)
	case "heartbeat":
		log.Printf("GDAX hearbeat: %s", msg.LastTradeId)
	case "recieved":
		writeGdaxOrder(msg)
	case "match":
		writeGdaxTrade(msg)
	case "open":
		writeGdaxBook(msg.OrderId, msg.Price, msg.RemainingSize, msg.Side, msg.Time)
	case "change":
		// price == "" means a change to a market order
		// (noise from gdax self-trade prevention system)
		if msg.Price != "" {
			//log.Printf("GDAX change: %s %s, %s", msg.OrderId, msg.Price, msg.NewSize)
			writeGdaxBook(msg.OrderId, msg.Price, msg.NewSize, msg.Side, msg.Time)
		}
	case "done":
		// price == "" indicates a filled market order
		// (duplicate of the "match" message)
		if msg.Price != "" {
			//log.Printf("GDAX done: %s %s, %s", msg.OrderId, msg.Price, msg.RemainingSize)
			// msg.RemainigSize will be > 0 in case of a cancelled order
			// (we already have this volume in the order book log)
			writeGdaxBook(msg.OrderId, msg.Price, "0", msg.Side, msg.Time)
		}
	}
}

func writeGdaxOrder(message GdaxMsg) {
}

func writeGdaxTrade(message GdaxMsg) {
}

func writeGdaxBook(orderId string, priceStr string, volumeStr string, side string, timeStr string) {
	price, err := strconv.ParseFloat(priceStr, 64)
	volume, err := strconv.ParseFloat(volumeStr, 64)
	timestamp, err := time.Parse(time.RFC3339Nano, timeStr) 
	if err != nil {
		log.Fatal("err: ", err)
	}
	queryStr := "INSERT INTO gdax_book_btcusd(order_id, price, volume, order_type, time_stamp, time_recieved) VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)"
	_, err = db.Exec(queryStr, orderId, price, volume, side, timestamp)
	if err != nil {
		log.Fatal("err: ", err)
	}
}


