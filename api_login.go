package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "login")
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
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
