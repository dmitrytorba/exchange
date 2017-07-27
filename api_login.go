package main

import (
	"github.com/dchest/captcha"
	"net/http"
)

const (
	TRIES_BEFORE_CAPTCHA = 2
	LOCKOUT              = 25
)

func loginHandler(w http.ResponseWriter, r *http.Request) error {
	count, err := checkLimit("login", r)
	if err != nil {
		return err
	}

	captcha := captcha.New()
	return executeTemplate(w, "login", 200, map[string]interface{}{
		"Captcha":   count >= TRIES_BEFORE_CAPTCHA,
		"CaptchaID": captcha,
	})
}

func loginPost(w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user := &User{
		Username: username,
		password: password,
	}

	err := authenticateByPassword(user)
	if err != nil {
		if err == ErrInvalidPassword || err == ErrUserNotFound {
			return executeTemplate(w, "login", 200, map[string]interface{}{
				"Username":  username,
				"Error":     "password was inccorect or username was not found",
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
