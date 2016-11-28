package main

import (
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) error {
	return executeTemplate(w, "login", 200, nil)
}

func loginPost(w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user := &User{
		username: username,
		password: password,
	}

	err := authenticateByPassword(user)
	if err != nil {
		if err == ErrInvalidPassword || err == ErrUserNotFound {
			return executeTemplate(w, "login", 200, map[string]interface{}{
				"Error":    "password or username was not found",
				"Username": username,
			})
		}

		return err
	}

	setCookie(user, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

// for use at the top of authenticated requests
func checkMe(r *http.Request) (*User, error) {
	return getUserFromCookie(r)
}

func logout(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return logoutFromCookie(r)
}
