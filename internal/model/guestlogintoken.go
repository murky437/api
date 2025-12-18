package model

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type GuestLoginToken struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func CreateGuestLoginToken(db *sql.DB, email string) (*GuestLoginToken, error) {
	guestLoginToken := &GuestLoginToken{
		Token: uuid.New().String(),
		Email: email,
	}
	_, err := db.Exec("INSERT INTO guest_login_token (token, email) VALUES (?, ?)", guestLoginToken.Token, guestLoginToken.Email)

	return guestLoginToken, err
}

func GetGuestLoginTokenByEmail(db *sql.DB, email string) (*GuestLoginToken, error) {
	guestLoginToken := &GuestLoginToken{
		Email: email,
	}
	err := db.QueryRow(`
		SELECT token
		FROM guest_login_token
		WHERE email = ?
	`, email).Scan(&guestLoginToken.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return guestLoginToken, nil
}

func GetOrCreateGuestLoginToken(db *sql.DB, email string) (*GuestLoginToken, error) {
	guestLoginToken, err := GetGuestLoginTokenByEmail(db, email)
	if err != nil {
		return nil, err
	}

	if guestLoginToken != nil {
		return guestLoginToken, nil
	}

	return CreateGuestLoginToken(db, email)
}

func GetGuestLoginTokenByToken(db *sql.DB, token string) (*GuestLoginToken, error) {
	guestLoginToken := &GuestLoginToken{
		Token: token,
	}
	err := db.QueryRow(`
		SELECT email
		FROM guest_login_token
		WHERE token = ?
	`, token).Scan(&guestLoginToken.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return guestLoginToken, nil
}

func DeleteGuestLoginTokenByToken(db *sql.DB, token string) error {
	_, err := db.Exec(`
		DELETE FROM guest_login_token
		WHERE token = ?
	`, token)

	return err
}
