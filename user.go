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
)

type User struct {
	id              int64
	username        string
	email           string
	password        string
	activationToken string
	sessionId       string
}

func decodeCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("dx45sp")
	if err != nil {
		return "", err
	}
	return url.QueryUnescape(cookie.Value)
}

func getUserFromCookie(r *http.Request) *User {
	sessionId, err := decodeCookie(r)
	if err != nil {
		return nil
	}
	return getSessionUser(sessionId)
}

func logoutFromCookie(r *http.Request) {
	sessionId, err := decodeCookie(r)
	if err != nil {
		return
	}
	removeSession(sessionId)
}

func addSession(user *User) {
	randBytes, err := scrypt.GenerateRandomBytes(32)
	if err != nil {
		panic(err)
	}
	sessionId := string(randBytes)
	// TODO: store more than a name
	err = rd.Set("session:"+sessionId, user.email, sessionTimeout).Err()
	if err != nil {
		panic(err)
	}
	user.sessionId = url.QueryEscape(sessionId)
}

func getSessionUser(sessionId string) *User {
	session, err := rd.Get("session:" + sessionId).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		panic(err)
	} else {
		var usr User
		// TODO: store more than a name
		usr.username = session
		return &usr
	}
}

func removeSession(sessionId string) {
	err := rd.Del("session:" + sessionId).Err()
	if err != nil {
		log.Printf("err removing session: %s", err)
	}
}

func findUserByEmail(email string) *User {
	var usr User
	queryStr := "select id, email from users where email = $1"
	err := db.QueryRow(queryStr, email).Scan(&usr.id, &usr.email)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		log.Fatal(err)
	default:
	}
	return &usr
}

func authenticateByToken(idStr string, token string) (*User, error) {
	usr := User{}
	var tokenHash string
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	queryStr := "select id, email, token from users where id = $1"
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

func authenticateByPassword(email string, password string) (*User, error) {
	var usr User
	var tokenHash string
	var passwordHash string
	queryStr := "select id, password, token from users where email = $1"
	err := db.QueryRow(queryStr, email).Scan(&usr.id, &passwordHash, &tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	err = scrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, err
	}
	usr.email = email
	addSession(&usr)
	return &usr, nil
}

func generateToken() string {
	randBytes, err := scrypt.GenerateRandomBytes(32)
	if err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(randBytes)
}

func encrypt(password string) []byte {
	hash, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}
	return hash
}

func createUser(email string, password string) (*User, error) {
	var usr User
	usr.email = email
	usr.activationToken = generateToken()
	tokenHash := encrypt(usr.activationToken)
	passwordHash := encrypt(password)

	queryStr := "INSERT INTO users(email, password, token) VALUES($1, $2, $3) returning id"
	err := db.QueryRow(queryStr, email, passwordHash, tokenHash).Scan(&usr.id)
	if err != nil {
		// check if the error is for a violation of a unique constraint like the username or email index
		if err.(*pq.Error).Code == "23505" { // 23505 is duplicate key value violates unique constraint
			switch err.(*pq.Error).Constraint {
			case "unique_username":
				return nil, ErrDuplicateUsername
			case "unique_email":
				return nil, ErrDuplicateEmail
			}
		}

		// all our other sql errors
		return nil, err
	}
	log.Printf("activation = %s", usr.activationToken)
	return &usr, nil
}
