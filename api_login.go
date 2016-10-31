package main

import (
	"fmt"
	"net/http"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) error {
	usr := User{}
	usr.username = r.FormValue("username")
	usr.password = r.FormValue("password")

	err := authenticateByPassword(&usr)
	if err != nil {
		return err
	}

	// no ID means no user was found
	if usr.id == 0 {
		w.WriteHeader(401)
		return nil
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
		return nil
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) error {
	usr := User{}
	usr.username = r.FormValue("username")
	usr.password = r.FormValue("password")
	usr.email = r.FormValue("email")
	return signup(&usr, w)
}

func signup(usr *User, w http.ResponseWriter) error {

	err := createUser(usr)

	if err != nil {
		switch err {
		case ErrDuplicateEmail:
			w.WriteHeader(400)
			fmt.Fprintln(w, "email already in use")
			return nil
		case ErrDuplicateUsername:
			w.WriteHeader(400)
			fmt.Fprintln(w, "username already in use")
			return nil
		default:
			return err
		}
	}

	fmt.Fprintln(w, usr.username)
	return nil
}

// for use at the top of authenticated requests
func checkMe(r *http.Request) (*User, error) {
	return getUserFromCookie(r)
}

func logout(w http.ResponseWriter, r *http.Request) error {
	return logoutFromCookie(r)
}
