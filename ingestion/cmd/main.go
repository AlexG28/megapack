package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"

	pb "github.com/AlexG28/megapack/proto/telemetry"
)

type TelemetryData struct {
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

func main() {
	time.Sleep(time.Second * 25)
	conn, err := openDBConnection()

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer conn.Close()

	if err = healthCheck(conn); err != nil {
		log.Fatalf("Health check failed: %v\n", err)
	}

	err = createTable(conn)

	if err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	}

	ch, q, err := openRabbitMQConnection("main")

	if err != nil {
		log.Fatalf("Rabbitmq error: %v", err)
	}

	fmt.Println("Successfully established all the major connections and ready to injest data into DB")

	defer ch.Close()

	var forever chan struct{}
	dataChan := make(chan TelemetryData, 100)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("consume error: %v\n", err)
	}

	var telData TelemetryData

	go func() {
		for d := range msgs {
			m := pb.TelemetryData{}
			if err := proto.Unmarshal(d.Body, &m); err != nil {
				log.Fatalf("error in decoding bytes into m")
				continue
			}

			telData = convertProtoToTelData(&m)
			dataChan <- telData
		}
	}()

	go addToDB(conn, dataChan)

	<-forever
}

func openDBConnection() (*pgx.Conn, error) {
	connStruct := pgx.ConnConfig{
		User:     "postgres",
		Password: "dbpassword",
		Host:     "timescaledb",
		Port:     5432,
		Database: "postgres",
	}

	conn, err := pgx.Connect(connStruct)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to timescaleDB: %v", err)
	}

	return conn, nil
}

func openRabbitMQConnection(queueName string) (*amqp.Channel, *amqp.Queue, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")

	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial rabbitmq: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create channel: %v", err)
	}
	// queueName = "main"
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	err = ch.Qos(
		1, 0, false,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set QoS: %v", err)
	}
	return ch, &q, nil
}

func healthCheck(conn *pgx.Conn) error {
	ctx := context.Background()

	err := conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("ingestion DB healthcheck failed: %v", err)
	}
	fmt.Println("Healthcheck Successfull!")
	return nil
}

func createTable(conn *pgx.Conn) error {
	var exists bool
	err := conn.QueryRow(`
	SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_name = 'telemetry_data'
	)`).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking if table exists: %w", err)
	}

	if exists {
		log.Println("The table already exists!")
		return nil
	}

	sql := `CREATE TABLE telemetry_data (
		unit_id VARCHAR(255),
		state VARCHAR(255),
		timestamp TIMESTAMPTZ,
		temperature FLOAT,
		charge INT,
		cycle INT,
		output INT,
		runtime INT,
		power INT);`

	_, err = conn.Exec(sql)

	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	fmt.Println("Created table")

	return nil
}

func addToDB(conn *pgx.Conn, dataChan <-chan TelemetryData) {
	for data := range dataChan {
		query := `INSERT INTO telemetry_data 
		(unit_id, state, timestamp, temperature, charge, cycle, output, runtime, power) 
		VALUES ($1, $2, $3::timestamptz, $4, $5, $6, $7, $8, $9)`

		_, err := conn.Exec(
			query,
			data.UnitID,
			data.State,
			data.Timestamp,
			data.TemperatureCelcius,
			data.ChargeLevelPercent,
			data.ChargeCycle,
			data.Output,
			data.Runtime,
			data.Power,
		)

		if err != nil {
			log.Printf("The error that occured in processData: %v\n", err)
		}

		if err != nil {
			log.Printf("%v\n", err)
		}
	}
}

func convertProtoToTelData(proto *pb.TelemetryData) TelemetryData {
	return TelemetryData{
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
