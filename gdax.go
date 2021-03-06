package main

import (
	"log"
	"encoding/json"
	"strconv"
	"time"
	"net/http"
	"io/ioutil"
	"gopkg.in/redis.v4"
	"strings"
)

const gdaxWS ="wss://ws-feed.gdax.com"
const gdaxREST = "https://api.gdax.com"

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
	MakerOrderId string `json:"maker_order_id"`
	TakerOrderId string `json:"taker_order_id"`
	NewSize string `json:"new_size"`
	OldSize string `json:"old_size"`
	NewFunds string
	OldFunds string
	LastTradeId string
	Message string
}

var gdaxSnapshotSequence map [string]int64
var gdaxEventBacklog map [string][]GdaxMsg

func connectGdax() {
	gdaxSnapshotSequence = make(map[string]int64)
	gdaxEventBacklog = make(map[string][]GdaxMsg)
	monitorGdaxSocket("btc", "usd")
	monitorGdaxSocket("eth", "usd")
	monitorGdaxSocket("eth", "btc")
}	

func monitorGdaxSocket(currencyBuy string, currencySell string) {
	currency := currencyBuy + currencySell
	currencyGdax := strings.ToUpper(currencyBuy) + "-" + strings.ToUpper(currencySell)
	gdaxSnapshotSequence[currency] = 0
	gdaxEventBacklog[currency] = make([]GdaxMsg, 0, 100)
	monitorWebsocket(
		gdaxWS,
		`{"type":"subscribe","product_ids":["` + currencyGdax + `"]}`,
		onGdaxEvent(currency))
	resetGdaxBook(currency, gdaxREST + "/products/" + currencyGdax + "/book?level=3")
}	

func onGdaxEvent(currency string) func(string) {
	return func(messageStr string) {
		var msg GdaxMsg
		json.Unmarshal([]byte(messageStr), &msg)
		//log.Printf("GDAX msg: %s", messageStr)

		if msg.Type == "error" {
			log.Printf("GDAX error: %s", msg.Message)
			return
		}
		if msg.Type == "heartbeat" {
			log.Printf("GDAX hearbeat: %s", msg.LastTradeId)
			return
		}
		if msg.Type == "received" {
			return
		}
		if msg.Type == "match" {
			writeGdaxTrade(msg, currency)
		}
		if msg.Price != "" {
			// when price == "" it means there is either a 'change' message for a market order
			// or a 'done' message that baiscally duplicates this order's 'match' message
			// (all this is noise that we can ignore)
			//writeGdaxBook(msg, currency)
		}
	}
}

type GdaxTrade struct {
	Price float64
	Volume float64
	Timestamp time.Time
	OrderId string
}

func parseGdaxTrade(msg GdaxMsg) GdaxTrade {
	var trade GdaxTrade
	var err error
	trade.Price, err = strconv.ParseFloat(msg.Price, 64)
	trade.Volume, err = strconv.ParseFloat(msg.Size, 64)
	trade.Timestamp, err = time.Parse(time.RFC3339Nano, msg.Time)
	if err != nil {
		log.Fatal("parse gdax trade err:", err)
	}
	trade.OrderId = msg.MakerOrderId
	return trade
}

func writeGdaxTrade(msg GdaxMsg, currency string) {
	trade := parseGdaxTrade(msg)
	queryStr := "INSERT INTO gdax_trades_" + currency +
		"(time_recieved,time_stamp,price,volume,order_id)" +
		"VALUES (CURRENT_TIMESTAMP,$1,$2,$3,$4)"
	_, err := db.Exec(queryStr, trade.Timestamp, trade.Price, trade.Volume, trade.OrderId)
	if err != nil {
		log.Fatal("trade insert err:", err)
	}
	jsonTrade, err := json.Marshal(trade)
	if err != nil {
		log.Fatal("trade to json marshal err:", err)
	}
	rd.Publish("gdax-trade-" + currency, string(jsonTrade))
}

func writeGdaxBook(msg GdaxMsg, currency string) {
	if gdaxSnapshotSequence[currency] == 0 {
		// no snapshot yet, save for later
		gdaxEventBacklog[currency] = append(gdaxEventBacklog[currency], msg)
	} else {
		price, err := strconv.ParseFloat(msg.Price, 64)
		timestamp, err := time.Parse(time.RFC3339Nano, msg.Time) 
		orderId := msg.OrderId
		entryType := "bid:" + msg.Type 
		if msg.Side == "sell" {
			entryType = "ask:" + msg.Type 
		}
		var volume float64
		if msg.Type == "match" {
			orderId = msg.MakerOrderId
			volume, err = strconv.ParseFloat(msg.Size, 64)
			volume = -volume
		} else if msg.Type == "change" {
			newVolume, err := strconv.ParseFloat(msg.NewSize, 64)
			oldVolume, err := strconv.ParseFloat(msg.OldSize, 64)
			if err != nil {
				log.Fatal("err: ", err)
			}
			volume = -(oldVolume-newVolume)
		} else {
			volume, err = strconv.ParseFloat(msg.RemainingSize, 64)
		}
		if err != nil {
			log.Fatal("writeGdaxBook err: ", err)
		}
		if msg.Type == "done" && msg.Reason == "canceled" {
			volume = -volume
		}
		queryStr := "INSERT INTO gdax_book_" + currency + "(order_id, price, volume, order_type, time_stamp, time_recieved) VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)"
		_, err = db.Exec(queryStr, orderId, price, volume, entryType, timestamp)
		if err != nil {
			log.Fatal("err: ", err)
		}
		key := currency + "-gdax-asks"
		if msg.Side == "buy" {
			key = currency + "-gdax-bids"
		}
		applyGdaxBookEntry(key, price, volume, orderId)
	}
}

type GdaxBookSnapshot struct {
	Sequence int64
	Bids [][]string
	Asks [][]string
}

type GdaxBookEntry struct {
	Price float64
	Volume float64
	OrderIds map[string]float64
}

func resetGdaxBook(currency string, url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var snapshot GdaxBookSnapshot
	json.Unmarshal(body, &snapshot)
	//log.Printf("seq: %s, bids: %s asks: %s", snapshot.Sequence, len(snapshot.Bids), len(snapshot.Asks))

	resetGdaxBookSide(currency, "asks", snapshot.Asks)
	resetGdaxBookSide(currency, "bids", snapshot.Bids)
	gdaxSnapshotSequence[currency] = snapshot.Sequence
	for _, msg := range gdaxEventBacklog[currency] {
		if msg.Sequence > gdaxSnapshotSequence[currency] {
			writeGdaxBook(msg, currency)
		}
	}
}

func resetGdaxBookSide(currency string, side string, msgs [][]string) {
	key := currency + "-gdax-" + side
	rd.Del(key)
	for _, msg := range msgs {
		price, err := strconv.ParseFloat(msg[0], 64)
		volume, err := strconv.ParseFloat(msg[1], 64)
		orderId := msg[2]
		if err != nil {
			log.Fatal(err)
		}
		entryType := side + ":snapshot"
		queryStr := "INSERT INTO gdax_book_" + currency + "(order_id, price, volume, order_type, time_recieved) VALUES($1, $2, $3, $4, CURRENT_TIMESTAMP)"
		_, err = db.Exec(queryStr, orderId, price, volume, entryType)
		if err != nil {
			log.Fatal("err: ", err)
		}
		applyGdaxBookEntry(key, price, volume, orderId)
	}
}

func applyGdaxBookEntry(key string, price float64, volume float64, orderId string) {
	priceStr := strconv.FormatFloat(price, 'f', -1, 64)
	vals, err := rd.ZRangeByScore(key, redis.ZRangeBy{
		Min: priceStr,
		Max: priceStr,
	}).Result()
	if err != nil || len(vals) > 1 {
		log.Fatal("applyGdaxBookEntry redis err:", err)
	}
	var entry GdaxBookEntry
	if len(vals) == 1 {
		json.Unmarshal([]byte(vals[0]), &entry)
		entry.Volume += volume
		if volume > 0 {
			entry.OrderIds[orderId] = volume
		} else if -(entry.OrderIds[orderId]) > volume {
			entry.OrderIds[orderId] += volume
		} else {
			delete(entry.OrderIds, orderId)
		}
		rd.ZRem(key, vals[0])
	} else {
		entry = GdaxBookEntry{
			Price: price,
			Volume: volume,
			OrderIds: map[string]float64 {
				orderId: volume,
			},
		}
	}

	if entry.Volume > 0.000001 {
		entryStr, err := json.Marshal(entry)
		if err != nil {
			log.Fatal("json parse error")
		}
		rd.ZAdd(key, redis.Z{
			Score: price,
			Member: entryStr,
		})
	}
}
