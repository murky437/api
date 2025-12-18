package model

import (
	"database/sql"
	"errors"
)

func LogDailyBrevoSend(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO
			daily_brevo_sends (number, date, updated_at)
		VALUES (1, STRFTIME('%Y-%m-%d','now'), STRFTIME('%Y-%m-%dT%H:%M:%fZ','now'))
		ON CONFLICT(date)
		DO UPDATE SET
		  number = number + 1,
		  updated_at = excluded.updated_at;
	`)

	return err
}

func GetDailyBrevoSendsNumber(db *sql.DB) (int64, error) {
	var number int64
	err := db.QueryRow(`
		SELECT
			number
		FROM
			daily_brevo_sends
		WHERE date = STRFTIME('%Y-%m-%d','now');
	`).Scan(&number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return number, nil
}
