package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	usr, err := authenticateByPassword(username, password)
	if err != nil {
		panic(err) // fix this later
	}

	if usr == nil {
		w.WriteHeader(401)
	} else {
		expire := time.Now().Add(10 * time.Minute)
		cookie := http.Cookie{
			Name:     "dx45sp",
			Value:    usr.sessionId,
			HttpOnly: true,
			Expires:  expire,
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
	user, err := checkMe(r)
	if err != nil {
		panic(err) // fix this later
	}

	if user == nil {
		fmt.Fprintf(w, "get out")
		return
	}

	fmt.Fprintf(w, "Hey %v!", user.username)
}

// for use at the top of authenticated requests
func checkMe(r *http.Request) (*User, error) {
	vars := mux.Vars(r)
	token := vars["token"]
	id := vars["id"]
	return authenticateByToken(id, token)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "logout")
}
