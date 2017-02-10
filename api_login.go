package main

import (
	"github.com/dchest/captcha"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) error {
	count, err := checkLimit("login", r)
	if err != nil {
		return err
	}

	captcha := captcha.New()
	return executeTemplate(w, "login", 200, map[string]interface{}{
		"Captcha":   count >= 1,
		"CaptchaID": captcha,
	})
}

func loginPost(w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// rate limit stuff
	count, err := rateLimit("login", r, 60*60)
	if err != nil {
		return err
	}
	if count > 25 { // should probably stop
		// process captcha
		return executeTemplate(w, "login", 200, map[string]interface{}{
			"Username": username,
			"Error":    "you have tried to log in too many times",
			"Captcha":  false,
		})
	}
	if count > 1 { // should probably check the captcha
		try := r.FormValue("captcha")
		id := r.FormValue("captchaID")

		if !captcha.VerifyString(id, try) { // captcha was wrong
			tryagain := captcha.New()
			return executeTemplate(w, "login", 200, map[string]interface{}{
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

	err = authenticateByPassword(user)
	if err != nil {
		if err == ErrInvalidPassword || err == ErrUserNotFound {
			return executeTemplate(w, "login", 200, map[string]interface{}{
				"Username": username,
				"Error":    "password was inccorect or username was not found",
				"Captcha":  count >= 3,
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
