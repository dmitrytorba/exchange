package main

import (
	"net/http"
	"log"
	"fmt"
	"io/ioutil"
)

func bitfinexTradesHandler(w http.ResponseWriter, r *http.Request) error {
	url := "https://api.bitfinex.com/v1/"
	req, err := http.NewRequest("GET", url+"/trades/btcusd", nil)
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

	fmt.Fprintf(w, string(body))
	return nil
}
