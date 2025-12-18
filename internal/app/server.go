package app

import (
	"fmt"
	"log"
	"net/http"
)

func StartApiServer(c *Container) {
	log.Printf("Starting API server at port %s...", c.Config.ApiPort)

	mux := NewMux(c)

	err := http.ListenAndServe(fmt.Sprintf(":%s", c.Config.ApiPort), mux)
	if err != nil {
		log.Println("Error starting API server:", err)
	}
}
