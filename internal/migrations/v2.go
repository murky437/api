package migrations

import "database/sql"

func v2(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE IF NOT EXISTS refresh_token (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	jwt TEXT NOT NULL,
	expires_at TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'NOW')),
	FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE
);
	`)

	return err
}
