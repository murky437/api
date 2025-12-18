package migrations

import "database/sql"

func v10(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE guest_login_token (
	token TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW'))
);
	`)

	return err
}
