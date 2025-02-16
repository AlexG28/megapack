package metrics

import (
	"fmt"
	"log"
	"monitoring/model"
)

func PrintFirstAndLast(arr []model.TelemetryData) {
	log.Println("Printing first and last:::")
	len := len(arr)
	fmt.Printf("arr[0]: %v\n", arr[0])
	fmt.Printf("arr[len-1]: %v\n", arr[len-1])
}
