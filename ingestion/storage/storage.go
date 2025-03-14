package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/AlexG28/megapack/ingestion/models"

	"github.com/jackc/pgx"
)

func AddToDB(conn *pgx.Conn, dataChan <-chan models.TelemetryData) {
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

func HealthCheck(conn *pgx.Conn) error {
	ctx := context.Background()
	err := conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("ingestion DB healthcheck failed: %v", err)
	}
	fmt.Println("Healthcheck Successfull!")
	return nil
}

func CreateTable(conn *pgx.Conn) error {
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

func Connect() (*pgx.Conn, error) {
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
