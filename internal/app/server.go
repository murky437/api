package app

import (
	"log"
	"net/http"
)

func StartApiServer(c *Container) {
	log.Printf("Starting API server at port 8080...")

	mux := NewMux(c)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println("Error starting API server:", err)
	}
}
