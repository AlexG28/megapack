package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type TelemetryData struct {
	UnitID             string  `json:"unit_id"`
	State              string  `json:"state"`
	Timestamp          string  `json:"timestamp"`
	TemperatureCelcius float32 `json:"temperature_celsius"`
	ChargeLevelPercent float32 `json:"charge_level_percent"`
	ChargeCycle        int     `json:"charge_cycle"`
	CumulativePower    int     `json:"cumulative_power"`
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

	err := SendToIngestion(telData)

	if err != nil {
		log.Printf("Unable to send to ingestion: %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "telemetry data received and processed")
}

func SendToIngestion(telData TelemetryData) error {
	url := "http://data-ingestion:8080/ingest"

	jsonData, err := json.Marshal(telData)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Incorrect Status: %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	return nil
}
