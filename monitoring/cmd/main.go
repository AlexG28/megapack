package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlexG28/megapack/monitoring/metrics"
	"github.com/AlexG28/megapack/monitoring/storage"
	"github.com/jackc/pgx"
)

func main() {
	log.Print("sleeping for 20 seconds")
	time.Sleep(time.Second * 20)
	log.Print("waking up")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()

	conn, err := storage.EstablishConnection()

	if err != nil {
		log.Fatalf(err.Error())
	}

	if err = storage.HealthCheck(conn); err != nil {
		log.Printf("Health check failed: %v\n", err)
	}

	defer conn.Close()

	go monitor(conn, ctx)

	<-ctx.Done()
	log.Println("Monitoring is shutting down")
}

func monitor(conn *pgx.Conn, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := metrics.Monitor(conn)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			time.Sleep(time.Second * 2)
		}
	}
}
