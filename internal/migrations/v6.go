package migrations

import "database/sql"

func v6(tx *sql.Tx) error {
	_, err := tx.Exec(`
ALTER TABLE user RENAME TO user_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
ALTER TABLE refresh_token RENAME TO refresh_token_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
ALTER TABLE project RENAME TO project_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
ALTER TABLE user_project RENAME TO user_project_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
ALTER TABLE long_reminder RENAME TO long_reminder_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
DROP INDEX idx_user_project_user_id;
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
CREATE TABLE user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'NOW')),
	updated_at DATETIME
);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
INSERT INTO user (id, username, password, created_at)
SELECT id, username, password, created_at
FROM user_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
CREATE TABLE refresh_token (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER REFERENCES user(id) ON DELETE CASCADE,
	jwt TEXT NOT NULL,
	expires_at DATETIME NOT NULL,
	created_at DATETIME NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW'))
);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
INSERT INTO refresh_token (id, user_id, jwt, expires_at, created_at)
SELECT id, user_id, jwt, expires_at, created_at
FROM refresh_token_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
CREATE TABLE project (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER REFERENCES user(id) ON DELETE CASCADE,
	title TEXT NOT NULL,
	slug TEXT NOT NULL UNIQUE,
	notes TEXT NOT NULL DEFAULT '',
	created_at DATETIME NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')),
	updated_at DATETIME
);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
CREATE INDEX idx_project_user_id ON project(user_id);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
INSERT INTO project (id, title, slug, notes, created_at, updated_at, user_id)
SELECT p.id, p.title, p.slug, p.notes, p.created_at, p.updated_at, up.user_id
FROM project_old p
LEFT JOIN (
	SELECT project_id, MIN(user_id) AS user_id
	FROM user_project_old
	GROUP BY project_id
) up ON up.project_id = p.id;
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
	created_at DATETIME DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')),
	updated_at DATETIME,
	last_reminded_at DATETIME,
	marked_done_at DATETIME,
	is_enabled INTEGER DEFAULT 1
);
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
INSERT INTO long_reminder (
	id, title, interval_days, user_id, created_at, updated_at,
	last_reminded_at, marked_done_at, is_enabled
)
SELECT id, title, interval_days, user_id, created_at, updated_at,
	   last_reminded_at, marked_done_at, is_enabled
FROM long_reminder_old;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
DROP TABLE refresh_token_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
DROP TABLE user_project_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
DROP TABLE project_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
DROP TABLE long_reminder_old;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
DROP TABLE user_old;
	`)

	return err
}
