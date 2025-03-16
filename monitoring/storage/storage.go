package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx"
)

func EstablishConnection() (*pgx.Conn, error) {
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

func HealthCheck(conn *pgx.Conn) error {
	ctx := context.Background()

	err := conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("monitoring DB healthcheck failed: %v", err)
	}
	log.Printf("Healthcheck Successfull!")
	return nil
}
