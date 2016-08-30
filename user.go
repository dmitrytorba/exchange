package main

import (
	"log"
	_ "database/sql"
	_ "github.com/lib/pq"
)

type User struct {
	id int
	email string
	password string
	config string
}

func findUserByEmail(email string) (User) {
	var usr User
	queryStr := "select id, email, config from users where email = ?"
	err := db.QueryRow(queryStr, email).Scan(&usr.id, &usr.email, &usr.config)
	if err != nil {
		log.Fatal(err)
	}
	return usr
}

func findUserByToken(token string) {

}

func createUser(email string, password string) {
	stmt, err := db.Prepare("INSERT INTO users(email, password) VALUES(?)")
	if err != nil {
		log.Fatal(err)
		}
	res, err := stmt.Exec("Dolly")
	if err != nil {
		log.Fatal(err)
		}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	
}
