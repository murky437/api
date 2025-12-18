package guestdbclear

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

func Handle(guestDb *sql.DB) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		log.Println("Clearing guestDb...")

		_, err := guestDb.Exec("PRAGMA foreign_keys = OFF;")
		if err != nil {
			log.Println(err)
			return err
		}
		defer func() {
			_, err := guestDb.Exec("PRAGMA foreign_keys = ON;")
			if err != nil {
				log.Println(err)
			}
		}()

		tx, err := guestDb.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		rows, err := tx.Query(`
			SELECT
			    name
			FROM
			    sqlite_master
			WHERE
			    type='table'
			  	AND name NOT LIKE 'sqlite_%'
			  	AND name NOT IN ('version');
		`)
		if err != nil {
			log.Println(err)
			return err
		}
		defer rows.Close()

		var table string
		for rows.Next() {
			err = rows.Scan(&table)
			if err != nil {
				log.Println(err)
				return err
			}
			_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s;", table))
			if err != nil {
				log.Println(err)
				return err
			}
		}

		err = rows.Err()
		if err != nil {
			log.Println(err)
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return err
		}

		log.Println("guestDb cleared")
		return nil
	}
}
