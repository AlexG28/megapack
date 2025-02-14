package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx"
)

func main() {
	fmt.Println("hello there, starting 5 second sleep.")
	time.Sleep(time.Second * 5)
	fmt.Println("Waking up!!")

	conn, err := establishConnection()

	if err != nil {
		log.Fatalf(err.Error())
	}

	if err = healthCheck(conn); err != nil {
		log.Printf("Health check failed: %v\n", err)
	}

	defer conn.Close()

	for i := range 10 {
		fmt.Printf("Attempt %v\n", i)
		if err := getSomeData(conn); err != nil {
			time.Sleep(time.Second * 1)
		}
		time.Sleep(time.Second * 2)
	}

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
		return nil, fmt.Errorf("Monitoring failed to connect to DB: %v\n", err)
	}

	log.Println("Monitoring successfully connected to DB!")

	return conn, nil
}

func healthCheck(conn *pgx.Conn) error {
	ctx := context.Background()

	err := conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("Monitoring DB healthcheck failed: %v\n", err)
	}
	log.Printf("Healthcheck Successfull!")
	return nil
}

func getSomeData(conn *pgx.Conn) error {
	err := currentRowCount(conn)
	if err != nil {
		log.Fatalf("rip: %v", err)
	}

	query := `
        SELECT unit_id, temperature_celsius, voltage_volts, charge_level_percent 
        FROM telemetry_data 
        LIMIT 10;
    `

	rows, err := conn.Query(query)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	log.Println("Successfully got rows")
	rowCount := 0
	for rows.Next() {
		rowCount++
		var (
			unitID             string
			temperatureCelsius float64
			voltageVolts       float64
			chargeLevelPercent float64
		)
		err := rows.Scan(&unitID, &temperatureCelsius, &voltageVolts, &chargeLevelPercent)
		if err != nil {
			log.Printf("error scanning row: %v", err)
		}

		fmt.Printf("Row %d: Unit ID: %s Temperature: %.2f°C, Voltage: %.2fV, Charge Level: %.2f\n",
			rowCount, unitID, temperatureCelsius, voltageVolts, chargeLevelPercent)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error after scanning rows: %v", err)
	}

	if rowCount == 0 {
		fmt.Println("No rows returned from the query.")
	}

	return nil
}

func currentRowCount(conn *pgx.Conn) error {
	var count int

	err := conn.QueryRow("SELECT COUNT(*) FROM telemetry_data").Scan(&count)

	if err != nil {
		return fmt.Errorf("Error occurred when counting lines: %v\n", err)
	}

	log.Printf("The row count is: %v\n", count)
	return nil
}
