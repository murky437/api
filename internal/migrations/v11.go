package migrations

import "database/sql"

func v11(tx *sql.Tx) error {
	_, err := tx.Exec(`
ALTER TABLE
    user
ADD COLUMN
    is_guest BOOLEAN NOT NULL DEFAULT 0
	`)

	return err
}
