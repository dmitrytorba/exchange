package main

import (
	"net/http"
	"unicode/utf8"
)

func signupPost(w http.ResponseWriter, r *http.Request) error {

	username := r.FormValue("username")
	password := r.FormValue("password")
	data := map[string]interface{}{
		"Error":    "too many accounnts have been made by this ip address, please wait",
		"Username": username,
	}

	// rate limit stuff
	count, err := rateLimit("signup", r, "", 60*30)
	if count > 1 { // cut-off for robots
		data["Captcha"] = true
	} else if count > 2 { // the 50 cut-off is to account for public places using our site
		data["Error"] = "too many accounnts have been made by this ip address, please wait"
		return executeTemplate(w, "signup", 200, data)
	}

	if utf8.RuneCountInString(username) == 0 || utf8.RuneCountInString(username) > 32 {
		data["Error"] = "usernames must be between 0 and 32 characters"
		return executeTemplate(w, "signup", 200, data)
	}

	if password != r.FormValue("password2") {
		data["Error"] = "passwords do not match"
		return executeTemplate(w, "signup", 200, data)
	}

	if utf8.RuneCountInString(password) < 3 || utf8.RuneCountInString(password) > 512 {
		data["Error"] = "password needs to be between 3 and 512 characters"
		return executeTemplate(w, "signup", 200, data)
	}

	user := &User{
		username: username,
		password: password,
	}

	err = createUser(user)
	if err != nil {
		if err == ErrDuplicateUsername {
			data["Error"] = "username has been taken"
			return executeTemplate(w, "signup", 200, data)
		}

		return err
	}

	setCookie(user, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func signupHandler(w http.ResponseWriter, r *http.Request) error {
	count, err := checkLimit("signup", r, "")
	if err != nil {
		return err
	}

	return executeTemplate(w, "signup", 200, map[string]interface{}{
		"Captcha": count > 1,
	})
}
