package main

import (
	"database/sql"
	"encoding/base64"
	"github.com/elithrar/simple-scrypt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"time"
	"gopkg.in/redis.v4"
	"net/url"
	"net/http"
)

const sessionTimeout = 10*time.Minute

type User struct {
	id              int64
	username        string
	email           string
	password        string
	activationToken string
	sessionId       string
}

func getUserFromCookie(r *http.Request) *User {
	cookie, err := r.Cookie("dx45sp")
	if err != nil {
		return nil
	}
	sessionId, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil
	}
	return getSessionUser(sessionId)
}

func addSession(user *User) {
	randBytes, err := scrypt.GenerateRandomBytes(32)
	if err != nil {
		panic(err)
	}
	sessionId := string(randBytes)
	// TODO: store more than a name
	err = rd.Set("session:" + sessionId, user.email, sessionTimeout).Err()
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

func authenticateByToken(idStr string, token string) *User {
	usr := User{}
	var tokenHash string
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil
	}
	queryStr := "select id, email, token from users where id = $1"
	err = db.QueryRow(queryStr, id).Scan(&usr.id, &usr.email, &tokenHash)
	if err != nil {
		return nil
	}
	err = scrypt.CompareHashAndPassword([]byte(tokenHash), []byte(token))
	if err != nil {
		return nil
	}
	return &usr
}

func authenticateByPassword(email string, password string) (*User) {
	var usr User
	var tokenHash string
	var passwordHash string
	queryStr := "select id, password, token from users where email = $1"
	err := db.QueryRow(queryStr, email).Scan(&usr.id, &passwordHash, &tokenHash)
	if err != nil {
		return nil
	}
	err = scrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil
	}
	usr.email = email
	addSession(&usr)
	return &usr
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

func createUser(email string, password string) *User {
	var usr User
	usr.email = email
	usr.activationToken = generateToken()
	tokenHash := encrypt(usr.activationToken)
	passwordHash := encrypt(password)
	queryStr := "INSERT INTO users(email, password, token) VALUES($1, $2, $3) returning id"
	err := db.QueryRow(queryStr, email, passwordHash, tokenHash).Scan(&usr.id)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("activation = %s", usr.activationToken)
	return &usr
}
