package worker

import (
	"log"
	"murky_api/internal/config"
	"time"

	"github.com/hibiken/asynq"
)

func StartScheduler(conf *config.Config) {
	log.Println("Scheduling worker tasks...")

	loc, err := time.LoadLocation(conf.Timezone)
	if err != nil {
		log.Println(err)
		return
	}

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: conf.RedisAddress},
		&asynq.SchedulerOpts{
			Location: loc,
		},
	)

	_, err = scheduler.Register("0 3 * * *", asynq.NewTask(TypeDbBackup, nil), asynq.MaxRetry(10))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = scheduler.Register("0 12 * * *", asynq.NewTask(TypeLongReminderPush, nil), asynq.MaxRetry(3))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = scheduler.Register("0 * * * *", asynq.NewTask(TypeGuestDbClear, nil), asynq.MaxRetry(3))
	if err != nil {
		log.Println(err)
		return
	}

	err = scheduler.Run()
	if err != nil {
		log.Println(err)
	}
}
