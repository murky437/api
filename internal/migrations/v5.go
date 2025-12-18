package migrations

import "database/sql"

func v5(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP INDEX IF EXISTS version_unique;`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DROP TABLE IF EXISTS schema_migrations;`)

	return err
}
