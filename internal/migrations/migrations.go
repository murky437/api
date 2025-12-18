package migrations

import (
	"database/sql"
	"errors"
	"log"
)

var migrations = []func(tx *sql.Tx) error{
	v1,
	v2,
	v3,
	v4,
	v5,
	v6,
	v7,
	v8,
	v9,
	v10,
	v11,
	v12,
	v13,
}

func Run(db *sql.DB, logger *log.Logger) error {

	err := createVersionTable(db)
	if err != nil {
		return err
	}

	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return err
	}

	logger.Println("Current database version:", currentVersion)

	newVersion := currentVersion
	migrationsApplied := 0

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := currentVersion; i < len(migrations); i++ {
		newVersion = i + 1

		logger.Printf("Applying migration v%d", newVersion)

		err = migrations[i](tx)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`INSERT INTO version (number) values (?)`, newVersion)
		if err != nil {
			return err
		}

		migrationsApplied++
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	if migrationsApplied != 0 {
		if migrationsApplied == 1 {
			logger.Printf("Applied 1 new migration")
		} else {
			logger.Printf("Applied %d new migrations", migrationsApplied)
		}
		logger.Println("New database version:", newVersion)
	} else {
		logger.Println("No new migrations applied")
	}

	return nil
}

func createVersionTable(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS version (
	number INTEGER PRIMARY KEY,
	applied_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'NOW'))
);
	`)
	return err
}

func getCurrentVersion(db *sql.DB) (int, error) {
	row := db.QueryRow(`SELECT number FROM version ORDER BY number DESC LIMIT 1`)
	var version int64
	err := row.Scan(&version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return int(version), nil
}
