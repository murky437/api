package longreminderpush

import (
	"context"
	"database/sql"
	"log"
	"murky_api/internal/clock"
	"murky_api/internal/config"
	"murky_api/internal/firebase"
	"murky_api/internal/model"

	"github.com/hibiken/asynq"
)

func Handle(db *sql.DB, conf *config.Config, firebaseMessageService firebase.MessageService, cl clock.Clock) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {

		log.Println("Checking due long reminders from db")

		remindersDueToday, err := model.GetLongRemindersDueToday(db, conf.Timezone)
		if err != nil {
			log.Println(err)
			return err
		}

		for _, reminder := range remindersDueToday {
			err := firebaseMessageService.Send(reminder.Title)
			if err != nil {
				log.Println(err)
				return err
			}

			err = model.SetLongReminderLastRemindedAt(db, reminder.Id, cl.Now().UTC())
			if err != nil {
				log.Println(err)
				return err
			}
		}

		return nil
	}
}
