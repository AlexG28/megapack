package metrics

import (
	"fmt"
	"monitoring/model"
)

func RequestsPerSecond(arr []model.TelemetryData) {
	len := len(arr)
	first := arr[0].Timestamp
	last := arr[len-1].Timestamp
	diff := first.Sub(last)

	fmt.Printf("Time Difference: %v   total seconds: %v\n", diff, diff.Seconds())
	fmt.Printf("Requests per second: %.3f\n", float64(len)/diff.Seconds())
}
