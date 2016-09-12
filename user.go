package main

import (
	"database/sql"
	"encoding/base64"
	"github.com/elithrar/simple-scrypt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"errors"
)

type User struct {
	id              int64
	email           string
	password        string
	activationToken string
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
	var usr User
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

func authenticateByPassword(email string, password string) (*User, error) {
	var usr User
	var tokenHash string
	var passwordHash string
	queryStr := "select id, password, token from users where email = $1"
	err := db.QueryRow(queryStr, email).Scan(&usr.id, &passwordHash, &tokenHash)
	if err != nil {
		return nil, err
	}
	if tokenHash != "" {
		return nil, errors.New("pending activation")
	}
	err = scrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, err
	}
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
