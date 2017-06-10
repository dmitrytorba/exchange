package main

import(
	"fmt"
	"time"
)

func backfillBitfinexTrades() {
	fmt.Println("Checking bitfinex trades history...")
	oldest, newest := getTimeRange("bitfinex_trades_btcusd")
	if oldest != nil {
		fmt.Println("found bitfinex: " + oldest.String() + " to " + newest.String())
	} else {
		fmt.Println("no history")
		//TODO: dowload history CSV from http://api.bitcoincharts.com/v1/csv/ 
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
