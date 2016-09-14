package main

import (
	_ "database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func api() (err error) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/settings", settingsHandler)
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/verify/{id}/{token}", verify)
	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/order", orderHandler).Methods("POST") // creating buy/sell orders
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Fatal(http.ListenAndServe(":4200", router))

	return nil
}
