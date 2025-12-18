package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: -db <db_file>")
	}
	dbPath := os.Args[2]

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO user (username, password) VALUES (?, ?)", "user", passwordHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added test user: 'user:pass'")
}
