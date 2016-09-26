package main

import (
	"net/http"
)

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "settings", 200, map[string]interface{}{})
}
