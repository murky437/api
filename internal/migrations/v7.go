package migrations

import "database/sql"

func v7(tx *sql.Tx) error {
	_, err := tx.Exec(`
ALTER TABLE version RENAME TO version_old;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
CREATE TABLE version (
	number INTEGER PRIMARY KEY,
	applied_at DATETIME NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'NOW'))
);
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
INSERT INTO version (number, applied_at)
SELECT number, applied_at
FROM version_old;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
DROP TABLE version_old;
	`)

	return err
}
