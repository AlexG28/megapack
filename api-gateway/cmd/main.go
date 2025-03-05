package main

import (
	"fmt"
	"gateway/handlers"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() { // open connection here
	time.Sleep(time.Second * 25)
	conn, ch, err := setupRabbitMQ("amqp://guest:guest@rabbitmq:5672/")

	if err != nil {
		log.Fatalf("Rabbitmq connection error: %v", err)
	}

	defer conn.Close()
	defer ch.Close()

	http.HandleFunc("/telemetry", handlers.TelemetryHandler(ch))
	http.HandleFunc("/health", handlers.HealthCheck)

	port := "8080"
	fmt.Printf("API gateway starting up on port %v\n", port)

	err = http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}

	fmt.Println("Listening and serving!")
}

func setupRabbitMQ(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, nil, fmt.Errorf("failure when dialing rabbitmq: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failure opening channel: %v", err)
	}
	fmt.Println("successfully setup rabbitmq")
	return conn, ch, nil
}
