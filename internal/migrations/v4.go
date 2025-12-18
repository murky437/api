package migrations

import "database/sql"

func v4(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE IF NOT EXISTS long_reminder (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	interval_days INTEGER NOT NULL,
	user_id INTEGER NOT NULL REFERENCES user(id),
	created_at TEXT DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'NOW')),
	updated_at TEXT,
	last_reminded_at TEXT,
	marked_done_at TEXT,
	is_enabled INTEGER DEFAULT 1
);
	`)

	_, err = tx.Exec(`
CREATE INDEX IF NOT EXISTS idx_long_reminder_user_id ON long_reminder(user_id);
	`)

	return err
}
