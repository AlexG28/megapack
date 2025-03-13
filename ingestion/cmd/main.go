package main

import (
	"fmt"
	"log"

	"github.com/jackc/pgx"
	"google.golang.org/protobuf/proto"

	pb "github.com/AlexG28/megapack/proto/telemetry"
)

func main() {
	// time.Sleep(time.Second * 25)
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
