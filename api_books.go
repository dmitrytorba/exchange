package main

import (
	"net/http"
	// "fmt"
	// "strconv"
	// "log"
)

func bitfinexBooksHandler(w http.ResponseWriter, r *http.Request) error {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	
	pubsub, err := rd.Subscribe("bitfinex")
	if err != nil {
    panic(err)
	}
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}

		if msg.Payload == "spread-change" {
			// fmt.Fprintf(w, "data: { bid: %f, ask: %f }\n\n", bitfinexBid, bitfinexAsk)
			f.Flush()
		}
	}
	return nil
}

func gdaxBooksHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}
