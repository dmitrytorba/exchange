package main

import (
	"github.com/dchest/captcha"
	"net/http"
	"unicode/utf8"
)

func signupPost(w http.ResponseWriter, r *http.Request) error {

	username := r.FormValue("username")
	password := r.FormValue("password")
	data := map[string]interface{}{
		"Username": username,
	}

	// checkin basic username stuffs
	if utf8.RuneCountInString(username) == 0 || utf8.RuneCountInString(username) > 32 {
		return executeTemplate(w, "signup", 200, map[string]interface{}{
			"Username": username,
			"Error":    "usernames must be between 0 and 32 characters",
		})
	}
	if password != r.FormValue("password2") {
		return executeTemplate(w, "signup", 200, map[string]interface{}{
			"Username": username,
			"Error":    "passwords do not match",
		})
	}
	if utf8.RuneCountInString(password) < 3 || utf8.RuneCountInString(password) > 512 {
		return executeTemplate(w, "signup", 200, map[string]interface{}{
			"Username": username,
			"Error":    "passwords need to be between 3 and 512 characters",
		})
	}

	// rate limit stuff
	count, err := rateLimit("signup", r, 60*30)
	if err != nil {
		return err
	}
	if count > 50 { // the 50 cut-off is to account for public places using our site
		return executeTemplate(w, "signup", 200, map[string]interface{}{
			"Username": username,
			"Error":    "too many accounts have been made by this computer, please wait",
		})
	}
	if count > 5 || 1 == 1 { // cut-off for robots
		try := r.FormValue("captcha")
		id := r.FormValue("captchaID")

		if !captcha.VerifyString(id, try) { // captcha was wrong
			tryagain := captcha.New()
			return executeTemplate(w, "signup", 200, map[string]interface{}{
				"Username":  username,
				"Error":     "your captcha was not correct",
				"Captcha":   true,
				"CaptchaID": tryagain,
			})
		}
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
	count, err := checkLimit("signup", r)
	if err != nil {
		return err
	}

	return executeTemplate(w, "signup", 200, map[string]interface{}{
		"Captcha":   count >= 5 || 1 == 1,
		"CaptchaID": captcha.New(),
	})
}
