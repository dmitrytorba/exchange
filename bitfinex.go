package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type tradeAPI interface {
	pushOrder() error
}

type bitfinexAPI struct {
	key    string
	secret string
}

// creates a market order using the key and secret you hopefully provided
// only market orders and btcusd for now
func (b *bitfinexAPI) marketOrder(amount, price int) error {

	req, err := http.NewRequest("POST", "https://api.bitfinex.com/v1/order/new", nil)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"request": "/v1/order/new",
		"nonce":   fmt.Sprintf("%v", time.Now().Unix()*10000),
		"symbol":  "btcusd",
		"amount":  amount,
		"price":   price,
		"type":    "market",
	}

	// if only there was a convenient protocol that could store information in a
	// packet of sorts, ill just json encode the information and call
	// it a day
	payload_json, _ := json.Marshal(payload)
	payload_enc := base64.StdEncoding.EncodeToString(payload_json)

	// how about we also make em run through an encryption scheme???
	sig := hmac.New(sha512.New384, []byte(b.secret))
	sig.Write([]byte(payload_enc))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-BFX-APIKEY", b.key)
	req.Header.Add("X-BFX-PAYLOAD", payload_enc)
	req.Header.Add("X-BFX-SIGNATURE", hex.EncodeToString(sig.Sum(nil)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}

func connectBitfinex() {
	monitorWebsocket("book", "BTCUSD", onBookMessage)
	monitorWebsocket("trades", "tBTCUSD", onTradeMessage)
}

func onTradeMessage(message string) {
	price, timestamp, volume := parseBitfinexTradeEntry(message)
	if price > 0 {
		writeTradeEntry(price, timestamp, volume)
	}
}

// bitfinex trade stream format:
// "[channel_id_int, event_string, [id_int, milli_time_int, volume_float, price_float]]"
// positive vol means 'buy', neg vol means 'sell'
// (meaningless to specify buy/sell for a trade, prob an artifact from the orderbook)
// there are two events for a trade, first a low-latency event_string='te'
// then a confirmation event_string='tu'
func parseBitfinexTradeEntry(entry string) (float64, time.Time, float64) {
	log.Printf(entry)
	entry = strings.Replace(entry, "[", "", -1)
	entry = strings.Replace(entry, "]", "", -1)
	parts := strings.Split(entry, ",")
	if len(parts) != 6 {
		// TODO support heartbeat?
		log.Printf("dont understand: %s", parts)
		return 0, time.Now(), 0
	} else {
		if parts[1] == `"te"` {
			price, err := strconv.ParseFloat(parts[5], 64)
			unixtime, err := strconv.ParseInt(parts[3], 10, 64)
			timestamp := time.Unix(unixtime/1000, unixtime%1000)
			volume, err := strconv.ParseFloat(parts[4], 64)
			if err != nil {
				log.Fatal(err)
			}
			// log.Printf("price: %s, count: %s, vol: %s", price, orderCount, volume)
			return price, timestamp, volume
		} else {
			// TODO: handle confirmation?
			return 0, time.Now(), 0
		}
	}
}

func writeTradeEntry(price float64, timestamp time.Time, volume float64) {
	log.Printf("price: %s, time: %s, vol: %s", price, timestamp, volume)
	queryStr := "INSERT INTO bitfinex_trades_btcusd(price, volume, time_stamp, time_recieved) VALUES($1, $2, $3, CURRENT_TIMESTAMP);"
	_, err := db.Exec(queryStr, price, volume, timestamp)
	if err != nil {
		// we are inserting a trade that already exists (same timestamp)
		return
	}
}

func onBookMessage(message string) {
	price, orderCount, volume := parseBitfinexBookEntry(message)
	if price > 0 {
		writeBookEntry(price, orderCount, volume)
	}
}

type handlerFunction func(string)

func monitorWebsocket(channel string, pair string, handler handlerFunction) {
	socket, _, err := websocket.DefaultDialer.Dial("wss://api2.bitfinex.com:3000/ws/2", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	payload := `{"event": "subscribe", "channel": "` + channel + `", "pair": "` + pair + `"}`
	log.Println("payload: ", payload)
	err = socket.WriteMessage(websocket.TextMessage, []byte(payload))
	if err != nil {
		log.Println("write:", err)
		return
	}

	go func() {
		for {
			_, message, err := socket.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			handler(string(message))
		}
	}()
}

// bitfinex book stream format:
// "[channel_id_int,[price_float,count_int,volume_float]]"
func parseBitfinexBookEntry(entry string) (float64, int64, float64) {
	entry = strings.Replace(entry, "[", "", -1)
	entry = strings.Replace(entry, "]", "", -1)
	parts := strings.Split(entry, ",")
	if len(parts) != 4 {
		// TODO support heartbeat
		log.Printf("dont understand: %s", parts)
		return 0, 0, 0
	} else {
		price, err := strconv.ParseFloat(parts[1], 64)
		orderCount, err := strconv.ParseInt(parts[2], 10, 64)
		volume, err := strconv.ParseFloat(parts[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		// log.Printf("price: %s, count: %s, vol: %s", price, orderCount, volume)
		return price, orderCount, volume
	}
}

func writeBookEntry(price float64, orderCount int64, volume float64) {
	orderType := "buy"
	if volume < 0 {
		// this is an 'ask' order
		orderType = "sell"
		volume *= -1
	}
	queryStr := "INSERT INTO bitfinex_book_btcusd(price, order_count, volume, order_type, time_stamp) VALUES($1, $2, $3, $4, CURRENT_TIMESTAMP)"
	_, err := db.Exec(queryStr, price, orderCount, volume, orderType)
	if err != nil {
		log.Fatal("insert err", err)
	}
}
