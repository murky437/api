package migrations

import "database/sql"

func v12(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE daily_brevo_sends (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    number INTEGER NOT NULL,
    date TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')),
    updated_at TEXT NOT NULL
)
	`)

	return err
}
