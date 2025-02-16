package model

import "time"

type TelemetryData struct {
	UnitID             string
	Timestamp          time.Time
	TemperatureCelcius float32
	VoltageVolts       float32
	ChargeLevelPercent float32
}

var Layout = "2006-01-02 15:04:05.999999-07"
