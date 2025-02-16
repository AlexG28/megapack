package model

type TelemetryData struct {
	UnitID             string  `json:"unit_id"`
	Timestamp          string  `json:"timestamp"`
	TemperatureCelcius float32 `json:"temperature_celsius"`
	VoltageVolts       float32 `json:"voltage_volts"`
	ChargeLevelPercent float32 `json:"charge_level_percent"`
}
