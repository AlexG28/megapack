package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type TelemetryData struct {
	UnitID             string  `json:"unit_id"`
	Timestamp          string  `json:"timestamp"`
	TemperatureCelcius float32 `json:"temperature_celsius"`
	VoltageVolts       float32 `json:"voltage_volts"`
	ChargeLevelPercent float32 `json:"charge_level_percent"`
}

func TelemetryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var telData TelemetryData

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&telData); err != nil {
		http.Error(w, "Error decoding Json", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	log.Printf("Received telemetry data: %+v\n", telData)

	// add validation here

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "telemetry data received and processed")
}
