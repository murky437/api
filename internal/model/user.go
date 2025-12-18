package model

import (
	"database/sql"
	"errors"
)

type User struct {
	Id       int
	Username string
	Password string
	IsGuest  bool
}

func FindUserByUsername(db *sql.DB, username string) (*User, error) {
	var user User
	query := "SELECT id, username, password, is_guest FROM user WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password, &user.IsGuest)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, err
}

func CreateGuestUser(db *sql.DB, username string) (*User, error) {
	user := &User{
		Username: username,
		Password: "",
		IsGuest:  true,
	}
	res, err := db.Exec(`
		INSERT INTO user (username, password, is_guest) VALUES (?, ?, ?)
	`, user.Username, user.Password, user.IsGuest)
	if err != nil {
		return nil, err
	}
	userId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.Id = int(userId)

	return user, nil
}

func UserExists(db *sql.DB, username string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)", username).Scan(&exists)
	return exists, err
}
