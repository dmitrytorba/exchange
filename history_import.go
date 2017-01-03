package main

import(
	"fmt"
	"time"
	"log"
	"strconv"
	"strings"
	"github.com/gorilla/websocket"
)

func loadHistory() {
	loadBitfinexTrades()
	logBitfinexBook()
}

func loadBitfinexTrades() {
	fmt.Println("Checking bitfinex trades history...")
	oldest, newest := getTimeRange("bitfinex_trades_btcusd")
	if oldest != nil {
		fmt.Println("found bitfinex: " + oldest.String() + " to " + newest.String())
	} else {
		fmt.Println("no history")
		//TODO: dowload CSV from http://api.bitcoincharts.com/v1/csv/ 
	}
}

func getTimeRange(table string) (oldest *time.Time, newest *time.Time) {
	queryStr := "select max(timestamp), min(timestamp) from $1"
	err := db.QueryRow(queryStr, table).Scan(newest, oldest)
	if err !=nil {
		return nil, nil
	}
	return oldest, newest
}

// bitfinex book stream format:
// "[channel_id_int,[price_float,count_int,volume_float]]"
func parseBitfinexBookEntry(entry string) {
	entry = strings.Replace(entry, "[", "", -1)
	entry = strings.Replace(entry, "]", "", -1)
	parts := strings.Split(entry, ",")
	if len(parts) != 4 {
		log.Printf("dont understand: %s", parts)
	} else {
		price, err := strconv.ParseFloat(parts[1], 64)
		orderCount, err := strconv.Atoi(parts[2])
		volume, err := strconv.ParseFloat(parts[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("price: %s, count: %s, vol: %s", price, orderCount, volume)
	}
}

func logBitfinexBook() {
	socket, _, err := websocket.DefaultDialer.Dial("wss://api2.bitfinex.com:3000/ws/2", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	err = socket.WriteMessage(websocket.TextMessage, []byte(`{"event": "subscribe", "channel": "Book", "pair": "BTCUSD"}`))
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
			parseBitfinexBookEntry(string(message))
		}
	}()
}
