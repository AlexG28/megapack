package main

import (
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

	testConnection(conn)

	fmt.Printf("seems to be working correctly???")
	os.Exit(0)
}

func testConnection(conn *pgx.Conn) {
	var greeting string
	err := conn.QueryRow("select 'Hello, Timescale!'").Scan(&greeting)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println(greeting)
}
