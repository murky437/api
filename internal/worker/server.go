package worker

import (
	"database/sql"
	"log"
	"murky_api/internal/clock"
	"murky_api/internal/config"
	"murky_api/internal/firebase"
	"murky_api/internal/s3"
	"murky_api/internal/worker/dbbackup"
	"murky_api/internal/worker/guestdbclear"
	"murky_api/internal/worker/longreminderpush"

	"github.com/hibiken/asynq"
)

func StartServer(db *sql.DB, conf *config.Config, s3Client s3.Client, firebaseMessageService firebase.MessageService, guestDb *sql.DB) {
	log.Println("Starting worker server...")

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: conf.RedisAddress},
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeDbBackup, dbbackup.Handle(db, conf, s3Client, clock.New()))
	mux.HandleFunc(TypeLongReminderPush, longreminderpush.Handle(db, conf, firebaseMessageService, clock.New()))
	mux.HandleFunc(TypeGuestDbClear, guestdbclear.Handle(guestDb))

	err := server.Run(mux)
	if err != nil {
		log.Println("Error starting asynq server:", err)
	}
}
