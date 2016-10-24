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
	emailToken      string
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
	// TODO: store more than the username
	err = rd.Set("session:"+sessionId, user.username, sessionTimeout).Err()
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
			return nil
		}
		return err
	}
	err = scrypt.CompareHashAndPassword([]byte(passwordHash), []byte(usr.password))
	if err != nil {
		return err
	}
	addSession(usr)
	return nil
}

func generateEmailToken() string {
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

func createUser(usr *User) error {
	passwordHash := encrypt(usr.password)
	usr.password = ""

 var err error
	if usr.email != "" {
		usr.emailToken = generateEmailToken()
		emailTokenHash := encrypt(usr.emailToken)

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
