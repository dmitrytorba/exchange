package main

import (
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "home", 200, map[string]interface{}{
		"Variable": "hello",
	})
}
