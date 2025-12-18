package migrations

import "database/sql"

func v1(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'NOW'))
);
	`)
	return err
}
