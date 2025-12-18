package migrations

import "database/sql"

func v8(tx *sql.Tx) error {
	_, err := tx.Exec(`
ALTER TABLE long_reminder RENAME TO long_reminder_old;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
DROP INDEX idx_long_reminder_user_id;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
CREATE TABLE long_reminder (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	interval_days INTEGER NOT NULL,
	user_id INTEGER NOT NULL REFERENCES user(id),
	created_at DATETIME NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')),
	updated_at DATETIME,
	last_reminded_at DATETIME,
	marked_done_at DATETIME,
	is_enabled BOOLEAN NOT NULL DEFAULT 1
);
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
INSERT INTO long_reminder (id, title, interval_days, user_id, created_at, updated_at, last_reminded_at, marked_done_at, is_enabled) 
SELECT id, title, interval_days, user_id, created_at, updated_at, last_reminded_at, marked_done_at, is_enabled
FROM long_reminder_old;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
CREATE INDEX idx_long_reminder_user_id ON long_reminder(user_id);
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
DROP TABLE long_reminder_old;
	`)

	return err
}
