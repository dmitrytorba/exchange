package main

import (
	_ "database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func signupPageHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "signup", 200, nil)
}

// route level error handling implemented at described in this article
// https://blog.golang.org/error-handling-and-go
type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		executeTemplate(w, "error", 500, map[string]interface{}{
			"Error": err.Error(),
		})

		// current error handling scheme is to just log to console
		// maybe one day add extra special high-level error handling/logging
		// for serious problems like db write failures on balances
		log.Println(err)
	}
}

func api() (err error) {
	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/", appHandler(homeHandler))
	router.Handle("/settings", appHandler(settingsHandler))
	router.Handle("/signup", appHandler(signupHandler)).Methods("POST")
	router.Handle("/login", appHandler(login))
	router.Handle("/logout", appHandler(logout))
	router.Handle("/order", appHandler(orderHandler)).Methods("POST") // creating buy/sell orders
	router.Handle("/history", appHandler(historyHandler))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Fatal(http.ListenAndServe(":4200", router))

	return nil
}
