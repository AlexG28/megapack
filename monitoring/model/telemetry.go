package model

import "time"

type TelemetryData struct {
	UnitID             string
	State              string
	Timestamp          time.Time
	TemperatureCelcius float32
	ChargeLevelPercent float32
	ChargeCycle        int
	CumulativePower    int
}

var Layout = "2006-01-02 15:04:05.999999-07"
