package main

import(
	"fmt"
	"time"
)

func loadHistory() {
	fmt.Println("Checking history...")
	loadBitfinex()
}

func loadBitfinex() {
	oldest, newest := getTimeRange("bitfinex")
	if oldest != nil {
		fmt.Println("found bitfinex: " + oldest.String() + " to " + newest.String())
	} else {
		fmt.Println("no history")
	}
}

func getTimeRange(exchange string) (oldest *time.Time, newest *time.Time) {
	queryStr := "select max(timestamp), min(timestamp) from history where exchange = $1"
	err := db.QueryRow(queryStr, exchange).Scan(newest, oldest)
	if err !=nil {
		return nil, nil
	}
	return oldest, newest
}
