package models

import (
	pb "github.com/AlexG28/megapack/proto/telemetry"
)

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

func ConvertProtoToTelData(proto *pb.TelemetryData) TelemetryData {
	return TelemetryData{
		UnitID:             proto.GetUnitId(),
		State:              proto.GetState(),
		Timestamp:          proto.GetTimestamp(),
		TemperatureCelcius: proto.GetTemperature(),
		ChargeLevelPercent: int(proto.GetCharge()),
		ChargeCycle:        int(proto.GetCycle()),
		Output:             int(proto.GetOutput()),
		Runtime:            int(proto.GetRuntime()),
		Power:              int(proto.GetPower()),
	}
}
