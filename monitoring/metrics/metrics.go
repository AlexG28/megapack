package metrics

import (
	"fmt"
	"time"

	"github.com/AlexG28/megapack/monitoring/model"
	"github.com/jackc/pgx"
)

func getCountByState(conn *pgx.Conn, states ...string) (int, error) {
	query := `SELECT COUNT(*) AS discharging_units
	FROM (
		SELECT DISTINCT ON (unit_id) unit_id, state
		FROM telemetry_data
		ORDER BY unit_id, timestamp::TIMESTAMP DESC
	) AS latest_entries
	WHERE state = ANY($1);`

	var count int
	err := conn.QueryRow(query, states).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("queryCount failed: %w", err)
	}
	return count, nil
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
		{&status.Fault, func() (int, error) { return getCountByState(conn, "fault", "maintenance") }},
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
