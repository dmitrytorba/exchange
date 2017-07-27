package main

import (
	"github.com/dchest/captcha"
	"net/http"
	"unicode/utf8"
)

const (
	SIGNUPS_BEFORE_CAPTCHA = 2
	//LOCKOUT                = 50
)

func signupPost(w http.ResponseWriter, r *http.Request) error {

	username := r.FormValue("username")
	password := r.FormValue("password")

	// rate limit stuff
	tryagain := captcha.New()
	count, err := rateLimit("signup", r, 60*30)
	if err != nil {
		return err
	}
	if count > LOCKOUT { // the 50 cut-off is to account for public places using our site
		return executeTemplate(w, "signup", 200, map[string]interface{}{
			"Username": username,
			"Error":    "too many accounts have been made by this computer, please wait",
		})
	}
	if count > SIGNUPS_BEFORE_CAPTCHA { // cut-off for robots
		try := r.FormValue("captcha")
		id := r.FormValue("captchaID")

		if !captcha.VerifyString(id, try) { // captcha was wrong
			return executeTemplate(w, "signup", 200, map[string]interface{}{
				"Username":  username,
				"Error":     "your captcha was not correct",
				"Captcha":   true,
				"CaptchaID": tryagain,
			})
		}
	}

	// checkin basic username stuffs
	if utf8.RuneCountInString(username) == 0 || utf8.RuneCountInString(username) > 32 {
		return executeTemplate(w, "signup", 200, map[string]interface{}{
			"Username":  username,
			"Error":     "usernames must be between 0 and 32 characters",
			"Captcha":   count > SIGNUPS_BEFORE_CAPTCHA,
			"CaptchaID": tryagain,
		})
	}
	if password != r.FormValue("password2") {
		return executeTemplate(w, "signup", 200, map[string]interface{}{
			"Username":  username,
			"Error":     "passwords do not match",
			"Captcha":   count > SIGNUPS_BEFORE_CAPTCHA,
			"CaptchaID": tryagain,
		})
	}
	if utf8.RuneCountInString(password) < 3 || utf8.RuneCountInString(password) > 512 {
		return executeTemplate(w, "signup", 200, map[string]interface{}{
			"Username":  username,
			"Error":     "passwords need to be between 3 and 512 characters",
			"Captcha":   count > SIGNUPS_BEFORE_CAPTCHA,
			"CaptchaID": tryagain,
		})
	}

	user := &User{
		Username: username,
		password: password,
	}

	err = createUser(user)
	if err != nil {
		if err == ErrDuplicateUsername {
			return executeTemplate(w, "signup", 200, map[string]interface{}{
				"Username": username,
				"Error":    "username has been taken",
			})
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
		"Captcha":   count >= SIGNUPS_BEFORE_CAPTCHA,
		"CaptchaID": captcha.New(),
	})
}
