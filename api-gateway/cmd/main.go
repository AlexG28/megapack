package main

import (
	"fmt"
	"gateway/handlers"
	"log"
	"net/http"
)

func main() {
	fmt.Println("hello there general kenobi")

	http.HandleFunc("/telemetry", handlers.TelemetryHandler)
	http.HandleFunc("/health", handlers.HealthCheck)

	port := "8080"
	log.Printf("API gateway starting up on port %v\n", port)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}

	log.Println("Listening and serving!")
}
