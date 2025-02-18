package metrics

import (
	"fmt"
	"monitoring/model"

	"github.com/jackc/pgx"
)

func RequestsPerSecond(arr []model.TelemetryData) {
	len := len(arr)
	first := arr[0].Timestamp
	last := arr[len-1].Timestamp
	diff := first.Sub(last)

	fmt.Printf("Requests per second: %.3f\n", float64(len)/diff.Seconds())
	fmt.Printf("first: %v\n", arr[0])
}

func AverageCharge(conn *pgx.Conn) error {
	query := `
        SELECT AVG(temperature_celsius) AS average_temperature
		FROM (
			SELECT temperature_celsius
			FROM telemetry_data
			ORDER BY timestamp DESC
			LIMIT 100
		) AS subquery;
    `
	rows, err := conn.Query(query)
	if err != nil {
		return fmt.Errorf("failed to perform query: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var out float32
		err := rows.Scan(&out)
		if err != nil {
			return fmt.Errorf("failed to convert from query to int: %v", err)
		}

		fmt.Printf("The out value is::::: %v\n", out)

	}

	return nil
}
