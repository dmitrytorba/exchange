package main

import (
	"log"
	"encoding/json"
	"strconv"
	"time"
	"net/http"
	"io/ioutil"
	"gopkg.in/redis.v4"
)

const gdaxWS ="wss://ws-feed.gdax.com"

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

var gdaxSnapshotSequence int64
var gdaxEventBacklog []GdaxMsg

func connectGdax() {
	gdaxSnapshotSequence = 0
	gdaxEventBacklog = make([]GdaxMsg, 0, 100)
	monitorWebsocket(
		gdaxWS,
		`{"type":"subscribe","product_ids":["BTC-USD"]}`,
		onGdaxEvent)
	resetGdaxBook()
}

func onGdaxEvent(messageStr string) {
	var msg GdaxMsg
	json.Unmarshal([]byte(messageStr), &msg)
	// log.Printf("GDAX msg: %s", messageStr)

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
		writeGdaxTrade(msg.MakerOrderId, msg.Price, msg.Size, msg.Time)
	}
	if msg.Price != "" {
		// when price == "" it means there is either a 'change' message for a market order
		// or a 'done' message that baiscally duplicates this order's 'match' message
		// (all this is noise that we can ignore)
		writeGdaxBook(msg)
	}
}

func writeGdaxTrade(orderId string, priceStr string, volumeStr string, timeStr string) {
	price, err := strconv.ParseFloat(priceStr, 64)
	volume, err := strconv.ParseFloat(volumeStr, 64)
	timestamp, err := time.Parse(time.RFC3339Nano, timeStr) 
	if err != nil {
		log.Fatal("err: ", err)
	}
	queryStr := "INSERT INTO gdax_trades_btcusd(order_id, price, volume, time_stamp, time_recieved) VALUES($1, $2, $3, $4, CURRENT_TIMESTAMP)"
	_, err = db.Exec(queryStr, orderId, price, volume, timestamp)
	if err != nil {
		log.Fatal("err: ", err)
	}
}

func writeGdaxBook(msg GdaxMsg) {
	if gdaxSnapshotSequence == 0 {
		// no snapshot yet, save for later
		gdaxEventBacklog = append(gdaxEventBacklog, msg)
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
		queryStr := "INSERT INTO gdax_book_btcusd(order_id, price, volume, order_type, time_stamp, time_recieved) VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)"
		_, err = db.Exec(queryStr, orderId, price, volume, entryType, timestamp)
		if err != nil {
			log.Fatal("err: ", err)
		}
		key := "gdax-asks"
		if msg.Side == "buy" {
			key = "gdax-bids"
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
	OrderIds []string
}

func resetGdaxBook() {
	url := "https://api.gdax.com"
	req, err := http.NewRequest("GET", url+"/products/BTC-USD/book?level=3", nil)
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
	log.Printf("seq: %s, bids: %s asks: %s", snapshot.Sequence, len(snapshot.Bids), len(snapshot.Asks))

	resetGdaxBookSide("gdax-asks", snapshot.Asks)
	resetGdaxBookSide("gdax-bids", snapshot.Bids)
	gdaxSnapshotSequence = snapshot.Sequence
	for _, msg := range gdaxEventBacklog {
		if msg.Sequence > gdaxSnapshotSequence {
			writeGdaxBook(msg)
		}
	}
}

func resetGdaxBookSide(key string, msgs [][]string) {
	rd.Del(key)
	for _, msg := range msgs {
		price, err := strconv.ParseFloat(msg[0], 64)
		volume, err := strconv.ParseFloat(msg[1], 64)
		orderId := msg[2]
		if err != nil {
			log.Fatal(err)
		}
		entryType := "bid:snapshot"
		if key == "gdax-asks" {
			entryType = "ask:snapshot"
		}
		queryStr := "INSERT INTO gdax_book_btcusd(order_id, price, volume, order_type, time_recieved) VALUES($1, $2, $3, $4, CURRENT_TIMESTAMP)"
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
		log.Fatal("redis err:", err)
	}
	var entry GdaxBookEntry
	if len(vals) == 1 {
		json.Unmarshal([]byte(vals[0]), &entry)
		entry.Volume += volume
		entry.OrderIds = append(entry.OrderIds, orderId)
		rd.ZRem(key, vals[0])
	} else {
		entry = GdaxBookEntry{
			Price: price,
			Volume: volume,
			OrderIds: []string{ orderId },
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
