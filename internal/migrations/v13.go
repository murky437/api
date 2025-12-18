package migrations

import "database/sql"

func v13(tx *sql.Tx) error {
	_, err := tx.Exec(`
ALTER TABLE project ADD COLUMN sort_index INTEGER NOT NULL DEFAULT 0;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
UPDATE
    project
SET sort_index = (
    SELECT count(*)
    FROM project p2
    WHERE p2.user_id = project.user_id
    	  AND p2.id < project.id
);
	`)

	return err
}
