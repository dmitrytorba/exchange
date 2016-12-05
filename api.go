package main

import (
	_ "database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v4"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"golang.org/x/net/http2"
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

// rateLimit is a function to be used by API paths to ensure no actor
// can aggressively use the API. It expects a name of the feature being rate limited (i.e. "signup")
// the request from which it will determine an ip address to limit, an optional username,
// and the expiration time in seconds.
func rateLimit(name string, r *http.Request, exp int) (int64, error) {
	ip := r.RemoteAddr
	if r.Header.Get("X-Forwarded-For") != "" {
		ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
		ip = ips[0]
	}

	key := fmt.Sprintf("%v:%v", name, ip)
	count, err := rd.Incr(key).Result()
	if err != nil {
		return 0, err
	}

	_, err = rd.Expire(key, time.Duration(exp)*time.Second).Result()
	if err != nil {
		return 0, err
	}

	fmt.Println(count)

	return count, nil
}

// checkLimit will return where the user is currently at on rate limits
func checkLimit(name string, r *http.Request) (int64, error) {
	ip := r.RemoteAddr
	if r.Header.Get("X-Forwarded-For") != "" { // I should double check that this actually gets the ip
		ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
		ip = ips[0]
	}

	key := fmt.Sprintf("%v:%v", name, ip)
	count, err := rd.Get(key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	number, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func api() (err error) {
	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/bitfinex/trades/btcusd", appHandler(bitfinexTradesHandler)).Methods("GET")

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

	server := &http.Server{
		Handler:      router,
	}
	http2.ConfigureServer(server, nil)
	log.Fatal(server.ListenAndServeTLS("localhost.cert", "localhost.key"))

	return nil
}
