package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx"
)

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

	healthCheck(conn)

	fmt.Print("Ready to accept data")

	queryCreateTable := `CREATE TABLE sensors (id SERIAL PRIMARY KEY, type VARCHAR(50), location VARCHAR(50));`
	_, err = conn.Exec(queryCreateTable)
	if err != nil {
		log.Fatalf("Unable to create SENSORS table: %v\n", err)
	}
	fmt.Println("Successfully created relational table SENSORS")

	os.Exit(0)
}

func healthCheck(conn *pgx.Conn) {
	ctx := context.Background()
	err := conn.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	log.Printf("Connection Successfull")
}
