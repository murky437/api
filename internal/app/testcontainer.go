package app

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"murky_api/internal/config"
	"murky_api/internal/constants"
	"murky_api/internal/jwt"
	"murky_api/internal/migrations"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func NewTestContainer(t *testing.T) *Container {
	conf := config.NewConfig()
	db := setupTestDatabase(t)
	jwtService := jwt.NewService(conf)
	return &Container{
		Config:     conf,
		Db:         db,
		JwtService: jwtService,
	}
}

func setupTestDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	err = setDatabasePragmas(db)
	if err != nil {
		t.Fatal(err)
	}

	err = migrations.Run(db, log.New(io.Discard, "", 0))
	if err != nil {
		t.Fatal(err)
	}

	insertInitialTestData(db, t)

	return db
}

func insertInitialTestData(db *sql.DB, t *testing.T) {
	// password = "pass"
	res, err := db.Exec(`INSERT INTO user (username, password) VALUES ('user', '$2a$10$U6X7NnivCljkduExJ9vqt.fqEGQrjxBczds1EbPrQjmLEw0eyUs9K')`)
	require.NoError(t, err)
	userId, err := res.LastInsertId()
	require.NoError(t, err)

	for i := 1; i <= 3; i++ {
		_, err := db.Exec(`INSERT INTO project (title, slug, user_id) VALUES (?, ?, ?)`, fmt.Sprintf("Project %d", i), fmt.Sprintf("project-%d", i), userId)
		require.NoError(t, err)
	}

	res, err = db.Exec(`INSERT INTO user (username, password) VALUES ('user2', '$2a$10$U6X7NnivCljkduExJ9vqt.fqEGQrjxBczds1EbPrQjmLEw0eyUs9K')`)
	require.NoError(t, err)
	user2Id, err := res.LastInsertId()
	require.NoError(t, err)

	for i := 4; i <= 6; i++ {
		_, err := db.Exec(`INSERT INTO project (title, slug, user_id) VALUES (?, ?, ?)`, fmt.Sprintf("Project %d", i), fmt.Sprintf("project-%d", i), user2Id)
		require.NoError(t, err)
	}

	insertReminderTestData(db, t, userId)
}

func insertReminderTestData(db *sql.DB, t *testing.T, userId int64) {
	updateTime, err := time.Parse(constants.SqliteDateTimeFormat, "2025-11-11T10:20:30.321Z")
	require.NoError(t, err)

	for i := 1; i <= 2; i++ {
		reminderId := fmt.Sprintf("longReminder%d", i)

		_, err := db.Exec(
			"INSERT INTO long_reminder (id, title, interval_days, user_id, updated_at) VALUES (?, ?, ?, ?, ?)",
			reminderId,
			fmt.Sprintf("Long reminder %d", i),
			40,
			userId,
			updateTime.Format(constants.SqliteDateTimeFormat),
		)
		require.NoError(t, err)
	}
}
