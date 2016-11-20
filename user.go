package main

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/elithrar/simple-scrypt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"gopkg.in/redis.v4"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const sessionTimeout = 10 * time.Minute

var (
	ErrDuplicateUsername = errors.New("username already exists in database")
	ErrDuplicateEmail    = errors.New("email already exists in the database")
	ErrUserNotFound      = errors.New("user was not found in the database")
)

type User struct {
	id         int64
	username   string
	email      string
	password   string
	emailToken string
	sessionId  string
}

func setCookie(usr *User, w http.ResponseWriter) {
	// set session cookie
	expire := time.Now().Add(10 * time.Minute)
	cookie := http.Cookie{
		Name:     "dx45sp",
		Value:    usr.sessionId,
		HttpOnly: true,
		Expires:  expire,
	}
	http.SetCookie(w, &cookie)
}
func decodeCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("dx45sp")
	if err != nil {
		return "", err
	}
	return url.QueryUnescape(cookie.Value)
}

func getUserFromCookie(r *http.Request) (*User, error) {
	sessionId, err := decodeCookie(r)
	if err != nil {
		return nil, nil
	}
	return getSessionUser(sessionId)
}

func logoutFromCookie(r *http.Request) error {
	sessionId, err := decodeCookie(r)
	if err != nil {
		return err
	}
	return removeSession(sessionId)
}

// addSession will save the sesion in redis
func addSession(user *User) error {
	randBytes, err := scrypt.GenerateRandomBytes(32)
	if err != nil {
		return err
	}
	sessionId := string(randBytes)
	// TODO: store more than the username
	err = rd.Set("session:"+sessionId, user.username, sessionTimeout).Err()
	if err != nil {
		return err
	}
	user.sessionId = url.QueryEscape(sessionId)
	return nil
}

// getSessionUser will retrieve the session from redis
func getSessionUser(sessionId string) (*User, error) {
	session, err := rd.Get("session:" + sessionId).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		var usr User
		// TODO: store more than a name
		usr.username = session
		return &usr, nil
	}
}

// removeSession will remove the session from redis
func removeSession(sessionId string) error {
	return rd.Del("session:" + sessionId).Err()
}

func findUserByEmail(email string) (*User, error) {
	var usr User
	queryStr := "select id, email from users where email = $1"
	err := db.QueryRow(queryStr, email).Scan(&usr.id, &usr.email)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	return &usr, nil
}

func authenticateByEmailToken(idStr string, token string) (*User, error) {
	usr := User{}
	var tokenHash string
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	queryStr := "select id, email, email_token from users where id = $1"
	err = db.QueryRow(queryStr, id).Scan(&usr.id, &usr.email, &tokenHash)
	if err != nil {
		// make sure we dont return an error for no users when we actually failed to find one
		if err == sql.ErrNoRows {
			return nil, nil
		}

		// normal SQL errors
		return nil, err
	}
	err = scrypt.CompareHashAndPassword([]byte(tokenHash), []byte(token))
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func authenticateByPassword(usr *User) error {
	var passwordHash string
	queryStr := "select id, password from users where username = $1"
	err := db.QueryRow(queryStr, usr.username).Scan(&usr.id, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return err
	}
	err = scrypt.CompareHashAndPassword([]byte(passwordHash), []byte(usr.password))
	if err != nil {
		return err
	}
	return addSession(usr)
}

func generateEmailToken() (string, error) {
	randBytes, err := scrypt.GenerateRandomBytes(32)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randBytes), nil
}

func encrypt(password string) ([]byte, error) {
	return scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
}

func createUser(usr *User) error {
	passwordHash, err := encrypt(usr.password)
	usr.password = ""
	if err != nil {
		return err
	}

	if usr.email != "" {
		usr.emailToken, err = generateEmailToken()
		if err != nil {
			return err
		}

		emailTokenHash, err := encrypt(usr.emailToken)
		if err != nil {
			return err
		}

		queryStr := "INSERT INTO users(username, password, email, email_token) VALUES($1, $2, $3, $4) returning id"
		err = db.QueryRow(queryStr, usr.username, passwordHash, usr.email, emailTokenHash).Scan(&usr.id)

	} else {
		queryStr := "INSERT INTO users(username, password) VALUES($1, $2) returning id"

		err = db.QueryRow(queryStr, usr.username, passwordHash).Scan(&usr.id)
	}

	if err != nil {
		// check if the error is for a violation of a unique constraint like the username or email index
		if err.(*pq.Error).Code == "23505" { // 23505 is duplicate key value violates unique constraint
			switch err.(*pq.Error).Constraint {
			case "unique_username":
				return ErrDuplicateUsername
			case "unique_email":
				return ErrDuplicateEmail
			}
		}

		// all our other sql errors
		return err
	}
	log.Printf("user %s created", usr.username)
	return nil
}
