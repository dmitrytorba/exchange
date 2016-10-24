package main

import (
	"fmt"
	"net/http"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {
	usr := User{}
	usr.username = r.FormValue("username")
	usr.password = r.FormValue("password")

	err := authenticateByPassword(&usr)
	if err != nil {
		panic(err) // TODO: fix this later
	}

	// no ID means no user was found
	if usr.id == 0 {
		w.WriteHeader(401)
	} else {
		// set session cookie
		expire := time.Now().Add(10 * time.Minute)
		cookie := http.Cookie{
			Name:     "dx45sp",
			Value:    usr.sessionId,
			HttpOnly: true,
			Expires:  expire,
		}
		http.SetCookie(w, &cookie)
		fmt.Fprintln(w, usr.username)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	usr := User{}
	usr.username = r.FormValue("username")
	usr.password = r.FormValue("password")
	usr.email = r.FormValue("email")
	signup(&usr, w)
}

func signup(usr *User, w http.ResponseWriter) {
	
	err := createUser(usr)

	if err != nil {
		switch err {
		case ErrDuplicateEmail:
			w.WriteHeader(400)
			fmt.Fprintln(w, "email already in use")
		case ErrDuplicateUsername:
			w.WriteHeader(400)
			fmt.Fprintln(w, "username already in use")
		default:
			panic(err)
		}
	}

	fmt.Fprintln(w, usr.username)
}

// for use at the top of authenticated requests
func checkMe(r *http.Request) *User {
	return getUserFromCookie(r)
}

func logout(w http.ResponseWriter, r *http.Request) {
	logoutFromCookie(r)
}
