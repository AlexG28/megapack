package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AlexG28/megapack/ingestion/models"

	"github.com/jackc/pgx"
)

func AddToDB(conn *pgx.Conn, data models.TelemetryData) error {
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
		return fmt.Errorf("the error that occured in processData: %v", err)
	}
	return nil
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

func ConnectToStorage() (*pgx.Conn, error) {
	numberOfRetries := 30
	var err error

	connStruct := pgx.ConnConfig{
		User:     "postgres",
		Password: "dbpassword",
		Host:     "timescaledb",
		Port:     5432,
		Database: "postgres",
	}

	for range numberOfRetries {
		conn, err := pgx.Connect(connStruct)

		if err == nil {
			return conn, nil
		}

		log.Println("failed to connect to timescaleDB, trying again")
		time.Sleep(2 * time.Second)

	}

	return nil, fmt.Errorf("failed to connect to timescaleDB: %v", err)
}
