package migrations

import "database/sql"

func v3(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE IF NOT EXISTS project (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	slug TEXT NOT NULL UNIQUE,
	notes TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'NOW')),
	updated_at TEXT
);
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
CREATE TABLE IF NOT EXISTS user_project (
	user_id    INTEGER NOT NULL,
	project_id INTEGER NOT NULL,
	PRIMARY KEY (user_id, project_id),
	FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
	FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE
);
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
CREATE INDEX IF NOT EXISTS idx_user_project_user_id ON user_project(user_id);
	`)

	return err
}
