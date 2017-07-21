package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	// "strconv"
	// "log"
)

func gdaxStatsHandler(w http.ResponseWriter, r *http.Request) error {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	vars := mux.Vars(r)
	currency := vars["currency"]
	key := "gdax-trade-" + currency
	pubsub, err := rd.Subscribe(key)
	if err != nil {
		panic(err)
	}
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "data: %s \n\n", msg.Payload)
		f.Flush()
	}
	return nil
}

func bitfinexBooksHandler( w http.ResponseWriter, r *http.Request) error {
	return nil
}
