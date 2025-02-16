package main

import (
	"context"
	"fmt"
	"log"
	"monitoring/metrics"
	"monitoring/model"
	"time"

	"github.com/jackc/pgx"
)

// type Row struct {
// 	unitID             string
// 	temperatureCelsius float64
// 	voltageVolts       float64
// 	chargeLevelPercent float64
// }

// type TelemetryData struct {
// 	UnitID             string  `json:"unit_id"`
// 	Timestamp          string  `json:"timestamp"`
// 	TemperatureCelcius float32 `json:"temperature_celsius"`
// 	VoltageVolts       float32 `json:"voltage_volts"`
// 	ChargeLevelPercent float32 `json:"charge_level_percent"`
// }

func main() {
	conn, err := establishConnection()
	log.Print("sleeping for 5 seconds")
	time.Sleep(time.Second * 5)
	log.Print("waking up")
	if err != nil {
		log.Fatalf(err.Error())
	}

	if err = healthCheck(conn); err != nil {
		log.Printf("Health check failed: %v\n", err)
	}

	defer conn.Close()

	last100Rows, err := getLast100Rows(conn)

	if err != nil {
		log.Fatalf("monitoring failed: %v\n", err)
	}

	metrics.PrintFirstAndLast(last100Rows)
}

func establishConnection() (*pgx.Conn, error) {
	connStruct := pgx.ConnConfig{
		User:     "postgres",
		Password: "teslagivemejob",
		Host:     "timescaledb",
		Port:     5432,
		Database: "postgres",
	}

	conn, err := pgx.Connect(connStruct)

	if err != nil {
		return nil, fmt.Errorf("monitoring failed to connect to DB: %v", err)
	}

	log.Println("Monitoring successfully connected to DB!")

	return conn, nil
}

func healthCheck(conn *pgx.Conn) error {
	ctx := context.Background()

	err := conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("monitoring DB healthcheck failed: %v", err)
	}
	log.Printf("Healthcheck Successfull!")
	return nil
}

func getLast100Rows(conn *pgx.Conn) ([]model.TelemetryData, error) {
	err := currentRowCount(conn)
	if err != nil {
		log.Fatalf("rip: %v", err)
	}

	query := `
        SELECT unit_id, timestamp, temperature_celsius, voltage_volts, charge_level_percent 
        FROM telemetry_data 
		ORDER BY timestamp DESC
        LIMIT 100;
    `

	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()
	rowCount := 0
	output := make([]model.TelemetryData, 100)

	for rows.Next() {
		var row model.TelemetryData
		err := rows.Scan(&row.UnitID, &row.Timestamp, &row.TemperatureCelcius, &row.VoltageVolts, &row.ChargeLevelPercent)
		if err != nil {
			log.Printf("error scanning row: %v", err)
		}
		output[rowCount] = row
		rowCount++
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %v", err)
	}

	if rowCount == 0 {
		fmt.Println("No rows returned from the query.")
	}

	return output, nil
}

func currentRowCount(conn *pgx.Conn) error {
	var count int

	err := conn.QueryRow("SELECT COUNT(*) FROM telemetry_data").Scan(&count)

	if err != nil {
		return fmt.Errorf("error occurred when counting lines: %v", err)
	}

	log.Printf("The row count is: %v\n", count)
	return nil
}
