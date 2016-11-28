package main

import (
	"net/http"
)

func signupPost(w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if password != r.FormValue("password2") {
		return executeTemplate(w, "signup", 200, map[string]interface{}{
			"Error":    "passwords do not match",
			"Username": username,
		})
	}

	user := &User{
		username: username,
		password: password,
	}

	err := createUser(user)
	if err != nil {
		if err == ErrDuplicateUsername {
			return executeTemplate(w, "signup", 200, map[string]interface{}{
				"Error":    "username has been taken",
				"Username": username,
			})
		}

		return err
	}

	setCookie(user, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func signupHandler(w http.ResponseWriter, r *http.Request) error {
	return executeTemplate(w, "signup", 200, nil)
}
