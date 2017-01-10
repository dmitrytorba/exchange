package main

import(
	"log"
	"strconv"
	"strings"
	"github.com/gorilla/websocket"
)

func loadHistory() {
	monitorWebsocket("book", "BTCUSD", onBookMessage)
	monitorWebsocket("trades", "tBTCUSD", onTradeMessage)
}

func onTradeMessage(message string) {
	log.Println(message)
}

func onBookMessage(message string) {
	price, orderCount, volume :=parseBitfinexBookEntry(message)
	if price != 0 {
		writeBookEntry(price, orderCount, volume)
	}
}

type handlerFunction func(string)

func monitorWebsocket(channel string, pair string, handler handlerFunction) {
	socket, _, err := websocket.DefaultDialer.Dial("wss://api2.bitfinex.com:3000/ws/2", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	payload :=`{"event": "subscribe", "channel": "` + channel + `", "pair": "` + pair + `"}`
	log.Println("payload: ", payload)
	err = socket.WriteMessage(websocket.TextMessage,  []byte(payload))
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
func parseBitfinexBookEntry(entry string) (float64, int, float64) {
	entry = strings.Replace(entry, "[", "", -1)
	entry = strings.Replace(entry, "]", "", -1)
	parts := strings.Split(entry, ",")
	if len(parts) != 4 {
		// TODO support heartbeat
		log.Printf("dont understand: %s", parts)
		return 0, 0, 0
	} else {
		price, err := strconv.ParseFloat(parts[1], 64)
		orderCount, err := strconv.Atoi(parts[2])
		volume, err := strconv.ParseFloat(parts[3], 64)
		if err != nil {	
			log.Fatal(err)
		}
		// log.Printf("price: %s, count: %s, vol: %s", price, orderCount, volume)
		return price, orderCount, volume
	}
}

func writeBookEntry(price float64, orderCount int, volume float64) {
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
