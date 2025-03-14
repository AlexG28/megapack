package main

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"

	"github.com/AlexG28/megapack/ingestion/message"
	"github.com/AlexG28/megapack/ingestion/models"
	"github.com/AlexG28/megapack/ingestion/storage"
	"github.com/jackc/pgx"
	"github.com/rabbitmq/amqp091-go"

	pb "github.com/AlexG28/megapack/proto/telemetry"
)

func main() {
	conn, err := storage.Connect()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	if err = storage.HealthCheck(conn); err != nil {
		log.Fatalf("Health check failed: %v\n", err)
	}

	if err := storage.CreateTable(conn); err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	}

	fmt.Println("Storage successfully connected to and ready to accept data.")

	ch, q, err := message.OpenRabbitMQConnection("main")

	if err != nil {
		log.Fatalf("Rabbitmq error: %v", err)
	}

	defer ch.Close()

	msgs, err := message.GetMessages(ch, q)

	if err != nil {
		log.Fatalf("consume error: %v\n", err)
	}

	fmt.Println("Successfully established all the major connections and ready to injest data into DB")

	var forever chan struct{}
	dataChan := make(chan models.TelemetryData, 100)
	go processMessages(msgs, dataChan)
	go storeTelemetry(conn, dataChan)

	<-forever
}

func storeTelemetry(conn *pgx.Conn, dataChan chan models.TelemetryData) {
	for {
		select {
		case data := <-dataChan:
			if err := storage.AddToDB(conn, data); err != nil {
				log.Printf("Failed to store telemetry: %v", err)
			}
		}
	}
}

func processMessages(msgs <-chan amqp091.Delivery, dataChan chan models.TelemetryData) {
	for {
		select {
		case msg := <-msgs:
			var tel pb.TelemetryData
			if err := proto.Unmarshal(msg.Body, &tel); err != nil {
				log.Printf("Failed to unmarshall message: %v", err)
				continue
			}

			dataChan <- models.ConvertProtoToTelData(&tel)
		}
	}
}
