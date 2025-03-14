package main

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"

	"github.com/AlexG28/megapack/ingestion/message"
	"github.com/AlexG28/megapack/ingestion/models"
	"github.com/AlexG28/megapack/ingestion/storage"
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
	go unMarshallMessages(msgs, dataChan)
	go storage.AddToDB(conn, dataChan)

	<-forever
}

func unMarshallMessages(msgs <-chan amqp091.Delivery, dataChan chan models.TelemetryData) {
	var telData models.TelemetryData
	for d := range msgs {
		m := pb.TelemetryData{}
		if err := proto.Unmarshal(d.Body, &m); err != nil {
			log.Fatalf("error in decoding bytes into m")
			continue
		}

		telData = convertProtoToTelData(&m)
		dataChan <- telData
	}
}

func convertProtoToTelData(proto *pb.TelemetryData) models.TelemetryData {
	return models.TelemetryData{
		UnitID:             proto.GetUnitId(),
		State:              proto.GetState(),
		Timestamp:          proto.GetTimestamp(),
		TemperatureCelcius: proto.GetTemperature(),
		ChargeLevelPercent: int(proto.GetCharge()),
		ChargeCycle:        int(proto.GetCycle()),
		Output:             int(proto.GetOutput()),
		Runtime:            int(proto.GetRuntime()),
		Power:              int(proto.GetPower()),
	}
}
