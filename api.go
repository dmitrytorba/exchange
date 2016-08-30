package main

import (
	"log"
	"net/http"
	_ "database/sql"
	"github.com/gorilla/mux"
)

func api() (err error) {
	router := mux.NewRouter().StrictSlash(true)
	
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/verify/{token}", verify)
	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)

	log.Fatal(http.ListenAndServe(":4200", router))

	return nil
}
