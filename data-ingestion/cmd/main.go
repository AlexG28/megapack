package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx"
)

type TelemetryData struct {
	UnitID             string  `json:"unit_id"`
	State              string  `json:"state"`
	Timestamp          string  `json:"timestamp"`
	TemperatureCelcius float32 `json:"temperature_celsius"`
	ChargeLevelPercent float32 `json:"charge_level_percent"`
	ChargeCycle        int     `json:"charge_cycle"`
	CumulativePower    int     `json:"cumulative_power"`
}

func main() {
	connStruct := pgx.ConnConfig{
		User:     "postgres",
		Password: "teslagivemejob",
		Host:     "timescaledb",
		Port:     5432,
		Database: "postgres",
	}

	conn, err := pgx.Connect(connStruct)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer conn.Close()

	if err = healthCheck(conn); err != nil {
		log.Printf("Health check failed: %v\n", err)
	}

	err = createTable(conn)

	if err != nil {
		log.Printf("Unable to create database: %v\n", err)
	}

	fmt.Print("Ready to accept data")

	dataChan := make(chan TelemetryData, 100)

	go processData(conn, dataChan)

	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var telData TelemetryData

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&telData); err != nil {
			http.Error(w, "Error decoding Json", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		dataChan <- telData

		w.WriteHeader(http.StatusOK)
	})
	port := "8080"
	log.Printf("API gateway starting up on port %v\n", port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
	log.Println("Listening and serving!")

}

func healthCheck(conn *pgx.Conn) error {
	ctx := context.Background()

	err := conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("Ingestion DB healthcheck failed: %v\n", err)
	}
	log.Printf("Healthcheck Successfull!")
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
		log.Printf("The table already exists!")
		return nil
	}

	sql := `CREATE TABLE telemetry_data (
		unit_id VARCHAR(255),
		state VARCHAR(255),
		timestamp VARCHAR(255),
		temperature_celsius FLOAT,
		charge_level_percent FLOAT,
		charge_cycle INT,
		cumulative_power INT
	);`

	_, err = conn.Exec(sql)

	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	log.Printf("Created table")

	return nil
}

func processData(conn *pgx.Conn, dataChan <-chan TelemetryData) {
	for data := range dataChan {
		query := `INSERT INTO telemetry_data (unit_id, state, timestamp, temperature_celsius, charge_level_percent, charge_cycle, cumulative_power) VALUES ($1, $2, $3::timestamptz, $4, $5, $6, $7)`

		_, err := conn.Exec(query, data.UnitID, data.State, data.Timestamp, data.TemperatureCelcius, data.ChargeLevelPercent, data.ChargeCycle, data.CumulativePower)

		if err != nil {
			log.Printf("The error that occured in processData: %v\n", err)
		}

		if err != nil {
			log.Printf("%v\n", err)
		}
	}
}
