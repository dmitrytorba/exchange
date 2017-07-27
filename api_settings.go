package main

import (
	"net/http"
)

func settingsPost(w http.ResponseWriter, r *http.Request) error {
	user, _ := getUserFromCookie(r)
	
	if user == nil {
		http.Redirect(w, r, "/login", 302)
		return nil 
	}

	user.gdaxKey = r.FormValue("gdax-key")
	user.gdaxSecret = r.FormValue("gdax-secret")
	user.gdaxPassphrase = r.FormValue("gdax-passpharse")

	updateUser(user)

	return nil
}

func settingsHandler(w http.ResponseWriter, r *http.Request) error {
	return executeTemplate(w, "settings", 200, map[string]interface{}{})
}
