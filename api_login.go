package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	usr := authenticateByPassword(username, password)

	if usr == nil {
		w.WriteHeader(401)
	} else {
		expire := time.Now().Add(10*time.Minute)
		cookie := http.Cookie{
			Name: "dx45sp",
			Value: usr.sessionId,
			HttpOnly: true,
			Expires: expire,
		}
		http.SetCookie(w, &cookie)  
		fmt.Fprintln(w, usr.email)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("username")
	password := r.FormValue("password")
	signup(email, password, w)
}

func signup(email string, password string, w http.ResponseWriter) {
	usr := findUserByEmail(email)

	if usr != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "email already in use")
	} else {
		usr = createUser(email, password)
	}
}

func verify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	id := vars["id"]
	authenticateByToken(id, token)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "logout")
}
