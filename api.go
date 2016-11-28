package main

import (
	_ "database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

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

	// order API
	router.Handle("/order", appHandler(orderHandler)).Methods("POST") // creating buy/sell orders

	// login API
	router.Handle("/login", appHandler(loginPost)).Methods("POST")
	router.Handle("/signup", appHandler(signupPost)).Methods("POST")
	router.Handle("/logout", appHandler(logout))

	// static pages
	router.Handle("/", appHandler(homeHandler))
	router.Handle("/login", appHandler(loginHandler))
	router.Handle("/signup", appHandler(signupHandler))
	router.Handle("/settings", appHandler(settingsHandler))
	router.Handle("/history", appHandler(historyHandler))

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Fatal(http.ListenAndServe(":4200", router))

	return nil
}
