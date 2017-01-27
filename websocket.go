package main

import (
	"log"
	"github.com/gorilla/websocket"
)

type handlerFunction func(string)

func monitorWebsocket(url string, payload string, handler handlerFunction) {
	socket, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
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
				log.Println(url + " disconnected: ", err)
				monitorWebsocket(url, payload, handler)
				return
			}
			handler(string(message))
		}
	}()
}
