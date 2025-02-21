package metrics

import (
	"fmt"
	"monitoring/model"
	"time"

	"github.com/jackc/pgx"
)

func queryCount(conn *pgx.Conn, query string, args ...interface{}) (int, error) {
	var count int
	err := conn.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("queryCount failed: %w", err)
	}
	return count, nil
}

func querySingleFloat(conn *pgx.Conn, query string, args ...interface{}) (float32, error) {
	var result float32
	err := conn.QueryRow(query, args...).Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("querySingleFloat failed: %w", err)
	}
	return result, nil
}

func getCountByState(conn *pgx.Conn, state string) (int, error) {
	return queryCount(conn,
		`SELECT COUNT(*) FROM telemetry_data WHERE state = $1`,
		state,
	)
}

func getCountByStates(conn *pgx.Conn, states ...string) (int, error) {
	query := `SELECT COUNT(*) FROM telemetry_data WHERE state = ANY($1)`
	return queryCount(conn, query, states)
}

func PerformMonitoring(conn *pgx.Conn) error {
	var status struct {
		Discharging int
		Charging    int
		Idle        int
		Fault       int
	}

	stateCounts := []struct {
		target   *int
		getCount func() (int, error)
	}{
		{&status.Discharging, func() (int, error) { return getCountByState(conn, "discharging") }},
		{&status.Charging, func() (int, error) { return getCountByState(conn, "charging") }},
		{&status.Idle, func() (int, error) { return getCountByState(conn, "idle") }},
		{&status.Fault, func() (int, error) { return getCountByStates(conn, "fault", "maintenance") }},
	}

	for _, sc := range stateCounts {
		count, err := sc.getCount()
		if err != nil {
			return fmt.Errorf("monitoring error: %w", err)
		}
		*sc.target = count
	}

	const (
		timeFormat   = "2006-01-02 15:04:05"
		separator    = "================================================"
		statusFormat = "%-19s  Charging: %-4d  Discharging: %-4d  Idle: %-4d  Faulty: %-4d\n"
	)

	timestamp := time.Now().UTC().Format(timeFormat)
	fmt.Println(separator)
	fmt.Printf(statusFormat, timestamp,
		status.Charging, status.Discharging, status.Idle, status.Fault)

	return nil
}

func RequestsPerSecond(arr []model.TelemetryData) {
	if len(arr) < 2 {
		fmt.Println("Insufficient data points for RPS calculation")
		return
	}

	duration := arr[0].Timestamp.Sub(arr[len(arr)-1].Timestamp).Seconds()
	if duration == 0 {
		fmt.Println("Zero duration between first and last timestamp")
		return
	}

	fmt.Printf("Requests per second: %.3f\n", float64(len(arr))/duration)
}

func averageCharge(conn *pgx.Conn) (float32, error) {
	const query = `
		SELECT AVG(temperature_celsius) 
		FROM (
			SELECT temperature_celsius
			FROM telemetry_data
			ORDER BY timestamp DESC
			LIMIT 100
		) AS subquery
	`
	return querySingleFloat(conn, query)
}
