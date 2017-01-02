package main

import(
	"fmt"
	"time"
	"log"
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
			log.Printf("recv: %s", message)
		}
	}()
}
