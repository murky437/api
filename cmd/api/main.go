package main

import (
	"log"
	"murky_api/internal/app"
	"murky_api/internal/worker"
)

func main() {
	c := app.NewContainer()
	defer func(c *app.Container) {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}(c)

	go worker.StartServer(c.Db, c.Config, c.S3Client, c.FirebaseMessageService, c.GuestDb)
	go worker.StartScheduler(c.Config)

	app.StartApiServer(c)
}
