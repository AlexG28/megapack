package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pb "github.com/AlexG28/megapack/proto/telemetry"
	"google.golang.org/protobuf/proto"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TelemetryDataStruct struct {
	UnitID             string  `json:"unit_id"`
	State              string  `json:"state"`
	Timestamp          string  `json:"timestamp"`
	TemperatureCelcius float32 `json:"temperature"`
	ChargeLevelPercent int     `json:"charge"`
	ChargeCycle        int     `json:"cycle"`
	Output             int     `json:"output"`
	Runtime            int     `json:"runtime"`
	Power              int     `json:"power"`
}

func TelemetryHandler(ch *amqp.Channel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var telData TelemetryDataStruct

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&telData); err != nil {
			http.Error(w, "Error decoding Json", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		err := sendToQueue(ch, "main", telData)

		if err != nil {
			log.Printf("Unable to send to ingestion: %v\n", err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "telemetry data received and processed")
	}
}

func sendToQueue(ch *amqp.Channel, queueName string, telData TelemetryDataStruct) error {
	telDataBytes := convertToProto(telData)
	protoData, err := proto.Marshal(telDataBytes)

	if err != nil {
		return fmt.Errorf("failed to convert to proto: %v", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        protoData,
		},
	)
	return err
}

func convertToProto(data TelemetryDataStruct) *pb.TelemetryData {
	return &pb.TelemetryData{
		UnitId:      data.UnitID,
		State:       data.State,
		Timestamp:   data.Timestamp,
		Temperature: data.TemperatureCelcius,
		Charge:      int32(data.ChargeLevelPercent),
		Cycle:       int32(data.ChargeCycle),
		Output:      int32(data.Output),
		Runtime:     int32(data.Runtime),
		Power:       int32(data.Power),
	}
}
