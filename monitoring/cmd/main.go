package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AlexG28/megapack/monitoring/metrics"
	"github.com/AlexG28/megapack/monitoring/model"

	"github.com/jackc/pgx"
)

func main() {
	conn, err := establishConnection()

	log.Print("sleeping for 2 seconds")
	time.Sleep(time.Second * 2)
	log.Print("waking up")

	if err != nil {
		log.Fatalf(err.Error())
	}

	if err = healthCheck(conn); err != nil {
		log.Printf("Health check failed: %v\n", err)
	}

	defer conn.Close()

	for range 50 {
		err := metrics.PerformMonitoring(conn)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		time.Sleep(time.Second * 2)
	}

}

func establishConnection() (*pgx.Conn, error) {
	connStruct := pgx.ConnConfig{
		User:     "postgres",
		Password: "dbpassword",
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
        SELECT 
			unit_id
			state,
			timestamp,
			temperature,
			charge,
			cycle,
			output,
			runtime,
			power
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
	var timestamp string

	for rows.Next() {
		var row model.TelemetryData
		err := rows.Scan(
			&row.UnitID,
			&row.State,
			&timestamp,
			&row.TemperatureCelcius,
			&row.ChargeLevelPercent,
			&row.ChargeCycle,
			&row.Output,
			&row.Runtime,
			&row.Power,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		parsedTime, err := time.Parse(model.Layout, timestamp)

		if err != nil {
			return nil, fmt.Errorf("error parsing time: %v", err)
		}

		row.Timestamp = parsedTime

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

	fmt.Printf("Row count: %v\n", count)

	return nil
}
