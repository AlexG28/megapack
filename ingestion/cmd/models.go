package main

type TelemetryData struct {
	UnitID             string  `json:"unit_id"`
	State              string  `json:"state"`
	Timestamp          string  `json:"timestamp"`
	TemperatureCelcius float32 `json:"temperature"`
	ChargeLevelPercent int     `json:"charge"`
	ChargeCycle        int     `json:"cycle"`
	Output             int     `json:"output"`
	Runtime            int     `json:"runtime"`
	Power              int     `json:"power"`
}
