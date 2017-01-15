package main

import (
	"log"
	"encoding/json"
	"strconv"
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
	OrderId string
	Reason string
	ProductId string
	Size float64
	Price string
	RemainingSize string `json:"remaining_size"`
	Sequence int64
	Time string
	OrderType string
	Funds string
	MakerOrderId string
	TakerOrderId string
	NewSize float64
	OldSize float64
	NewFunds float64
	OldFunds float64
	LastTradeId string
	Message string
}

func onGdaxEvent(messageStr string) {
	var message GdaxMsg
	json.Unmarshal([]byte(messageStr), &message)
	switch message.Type {
	case "error":
		log.Printf("GDAX error: %s", message.Message)
	case "heartbeat":
		log.Printf("GDAX hearbeat: %s", message.LastTradeId)
	case "recieved":
		writeGdaxOrder(message)
	case "match":
		writeGdaxTrade(message)
	default:
		writeGdaxBook(message)
	}
}

func writeGdaxOrder(message GdaxMsg) {
}

func writeGdaxTrade(message GdaxMsg) {
}

func writeGdaxBook(msg GdaxMsg) {
	var err error
	switch msg.Type {
	case "open":
		log.Printf("GDAX open: %s, %s", msg.Price, msg.RemainingSize)
		price, err := strconv.ParseFloat(msg.Price, 64)
		volume, err := strconv.ParseFloat(msg.RemainingSize, 64)
		if err != nil {
			log.Fatal("err: ", err)
		}
		queryStr := "INSERT INTO gdax_book_btcusd(price, volume, order_type, time_stamp) VALUES($1, $2, $3, CURRENT_TIMESTAMP)"
		_, err = db.Exec(queryStr, price, volume, msg.Side)
	case "change":
	case "done":
	}
	if err != nil {
		log.Fatal("err: ", err)
	}
}


