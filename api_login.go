package main

import (
	"fmt"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) error {
	usr := User{}
	usr.username = r.FormValue("username")
	usr.password = r.FormValue("password")

	// check the user's password
	err := authenticateByPassword(&usr)
	if err == ErrUserNotFound {
		w.WriteHeader(401)
		return nil
	}
	if err != nil {
		return err
	}

	setCookie(&usr, w)
	fmt.Fprintln(w, usr.username)
	return nil
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

	setCookie(usr, w)
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
