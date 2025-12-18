package model

import (
	"database/sql"
	"murky_api/internal/constants"
	"time"
)

type LongReminder struct {
	Id             string
	Title          string
	IntervalDays   int
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	LastRemindedAt *time.Time
	MarkedDoneAt   *time.Time
	IsEnabled      bool
}

func CreateLongReminder(db *sql.DB, reminder LongReminder, userId int) (LongReminder, error) {
	tx, err := db.Begin()
	if err != nil {
		return reminder, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO long_reminder (id, title, interval_days, user_id) VALUES (?, ?, ?, ?)", reminder.Id, reminder.Title, reminder.IntervalDays, userId)
	if err != nil {
		return reminder, err
	}

	err = tx.Commit()
	return reminder, err
}

func GetLongReminderByUserIdAndId(db *sql.DB, userId int, reminderId string) (LongReminder, error) {
	var r LongReminder

	var createdAt string
	var updatedAt, lastRemindedAt, markedDoneAt sql.NullString

	err := db.QueryRow(`
		SELECT r.id, r.title, r.interval_days, r.created_at, r.updated_at, r.last_reminded_at, r.marked_done_at, r.is_enabled
		FROM long_reminder r
		WHERE r.user_id = ? AND r.id = ?
	`, userId, reminderId).Scan(&r.Id, &r.Title, &r.IntervalDays, &createdAt, &updatedAt, &lastRemindedAt, &markedDoneAt, &r.IsEnabled)
	if err != nil {
		return r, err
	}

	r.CreatedAt, err = time.Parse(constants.SqliteDateTimeFormat, createdAt)
	if err != nil {
		return r, err
	}

	if updatedAt.Valid {
		t, err := time.Parse(constants.SqliteDateTimeFormat, updatedAt.String)
		if err != nil {
			return r, err
		}
		r.UpdatedAt = &t
	}
	if lastRemindedAt.Valid {
		t, err := time.Parse(constants.SqliteDateTimeFormat, lastRemindedAt.String)
		if err != nil {
			return r, err
		}
		r.LastRemindedAt = &t
	}
	if markedDoneAt.Valid {
		t, err := time.Parse(constants.SqliteDateTimeFormat, markedDoneAt.String)
		if err != nil {
			return r, err
		}
		r.MarkedDoneAt = &t
	}

	return r, nil
}

func GetLongRemindersByUserId(db *sql.DB, userId int) ([]LongReminder, error) {
	var reminders []LongReminder

	rows, err := db.Query(`
		SELECT r.id, r.title, r.interval_days, r.created_at, r.updated_at, r.last_reminded_at, r.marked_done_at, r.is_enabled
		FROM long_reminder r
		WHERE r.user_id = ?
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r LongReminder

		var createdAt string
		var updatedAt, lastRemindedAt, markedDoneAt sql.NullString

		if err := rows.Scan(&r.Id, &r.Title, &r.IntervalDays, &createdAt, &updatedAt, &lastRemindedAt, &markedDoneAt, &r.IsEnabled); err != nil {
			return nil, err
		}

		r.CreatedAt, err = time.Parse(constants.SqliteDateTimeFormat, createdAt)
		if err != nil {
			return nil, err
		}

		if updatedAt.Valid {
			t, err := time.Parse(constants.SqliteDateTimeFormat, updatedAt.String)
			if err != nil {
				return nil, err
			}
			r.UpdatedAt = &t
		}
		if lastRemindedAt.Valid {
			t, err := time.Parse(constants.SqliteDateTimeFormat, lastRemindedAt.String)
			if err != nil {
				return nil, err
			}
			r.LastRemindedAt = &t
		}
		if markedDoneAt.Valid {
			t, err := time.Parse(constants.SqliteDateTimeFormat, markedDoneAt.String)
			if err != nil {
				return nil, err
			}
			r.MarkedDoneAt = &t
		}

		reminders = append(reminders, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reminders, nil
}

func DeleteLongReminderByUserIdAndId(db *sql.DB, userId int, id string) (int64, error) {
	res, err := db.Exec(`
		DELETE FROM long_reminder
		WHERE user_id = ? AND id = ?
	`, userId, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func UpdateLongReminder(db *sql.DB, id string, reminder LongReminder) error {
	markedDoneAt := any(nil)
	if reminder.MarkedDoneAt != nil {
		markedDoneAt = reminder.MarkedDoneAt.Format(constants.SqliteDateTimeFormat)
	}

	_, err := db.Exec(`
		UPDATE long_reminder 
		SET title = :title, interval_days = :intervalDays, marked_done_at = :markedDoneAt, is_enabled = :isEnabled, updated_at = STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')
		WHERE id = :id
	`,
		sql.Named("title", reminder.Title),
		sql.Named("intervalDays", reminder.IntervalDays),
		sql.Named("markedDoneAt", markedDoneAt),
		sql.Named("isEnabled", reminder.IsEnabled),
		sql.Named("id", id),
	)
	return err
}

func GetLongRemindersDueToday(db *sql.DB, timezone string) ([]LongReminder, error) {
	var reminders []LongReminder

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}
	year, month, day := time.Now().In(loc).Date()
	todayUtc := time.Date(year, month, day, 0, 0, 0, 0, loc).UTC()
	tomorrowUtc := todayUtc.AddDate(0, 0, 1)

	// TODO: fix this after writing a test, using SELECT alias in WHERE isn't technically valid sql (but sqlite allows it)
	rows, err := db.Query(`
		SELECT
		    r.id, r.title, r.interval_days, r.created_at, r.updated_at, r.last_reminded_at, r.marked_done_at, r.is_enabled,
			IIF(
				r.marked_done_at IS NULL,
				r.created_at,
				STRFTIME('%Y-%m-%dT%H:%M:%fZ', r.marked_done_at, CONCAT('+', r.interval_days, ' days'))
			) AS due_date
		FROM long_reminder r
		JOIN user u ON r.user_id = u.id
		WHERE
		    u.is_guest = false
		    AND r.is_enabled = true
		  	AND due_date < :tomorrowDate
		  	AND (
		  	    r.last_reminded_at IS NULL
		  	    OR (
		  	        due_date >= :todayDate
		  	        AND r.last_reminded_at < :todayDate
		  	    )
		  	    OR (
		  	        due_date < :todayDate
		  	        	-- TODO: rethink the query, this part is wrong, this will repeat yesterdays reminder 1 time next day
		  	        AND r.last_reminded_at < due_date
		  	    )
		  	);
	`,
		sql.Named("todayDate", todayUtc.Format(constants.SqliteDateTimeFormat)),
		sql.Named("tomorrowDate", tomorrowUtc.Format(constants.SqliteDateTimeFormat)),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r LongReminder

		var createdAt string
		var updatedAt, lastRemindedAt, markedDoneAt sql.NullString
		var dueDate string

		if err := rows.Scan(&r.Id, &r.Title, &r.IntervalDays, &createdAt, &updatedAt, &lastRemindedAt, &markedDoneAt, &r.IsEnabled, &dueDate); err != nil {
			return nil, err
		}

		r.CreatedAt, err = time.Parse(constants.SqliteDateTimeFormat, createdAt)
		if err != nil {
			return nil, err
		}

		if updatedAt.Valid {
			t, err := time.Parse(constants.SqliteDateTimeFormat, updatedAt.String)
			if err != nil {
				return nil, err
			}
			r.UpdatedAt = &t
		}
		if lastRemindedAt.Valid {
			t, err := time.Parse(constants.SqliteDateTimeFormat, lastRemindedAt.String)
			if err != nil {
				return nil, err
			}
			r.LastRemindedAt = &t
		}
		if markedDoneAt.Valid {
			t, err := time.Parse(constants.SqliteDateTimeFormat, markedDoneAt.String)
			if err != nil {
				return nil, err
			}
			r.MarkedDoneAt = &t
		}

		reminders = append(reminders, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reminders, nil
}

func SetLongReminderLastRemindedAt(db *sql.DB, id string, lastRemindedAt time.Time) error {
	_, err := db.Exec(`
		UPDATE long_reminder 
		SET last_reminded_at = :lastRemindedAt
		WHERE id = :id
	`,
		sql.Named("lastRemindedAt", lastRemindedAt.Format(constants.SqliteDateTimeFormat)),
		sql.Named("id", id),
	)
	return err
}
