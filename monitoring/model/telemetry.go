package model

import "time"

type TelemetryData struct {
	UnitID             string
	State              string
	Timestamp          time.Time
	TemperatureCelcius float32
	ChargeLevelPercent int
	ChargeCycle        int
	Output             int
	Runtime            int
	Power              int
}

var Layout = "2006-01-02 15:04:05.999999-07"
