package dirtrally

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"time"
)

type TelemetryFrame struct {
	Time                float32
	LapTime             float32
	LapDistance         float32
	TotalDistance       float32
	X                   float32 // World space position
	Y                   float32 // World space position
	Z                   float32 // World space position
	Speed               float32
	Xv                  float32 // Velocity in world space
	Yv                  float32 // Velocity in world space
	Zv                  float32 // Velocity in world space
	Xr                  float32 // World space right direction
	Yr                  float32 // World space right direction
	Zr                  float32 // World space right direction
	Xd                  float32 // World space forward direction
	Yd                  float32 // World space forward direction
	Zd                  float32 // World space forward direction
	SuspPosBL           float32
	SuspPosBR           float32
	SuspPosFL           float32
	SuspPosFR           float32
	SuspVelBL           float32
	SuspVelBR           float32
	SuspVelFL           float32
	SuspVelFR           float32
	WheelSpeedBL        float32
	WheelSpeedBR        float32
	WheelSpeedFL        float32
	WheelSpeedFR        float32
	Throttle            float32
	Steer               float32
	Brake               float32
	Clutch              float32
	Gear                float32
	GForceLat           float32
	GForceLon           float32
	Lap                 float32
	EngineRate          float32
	SliProNativeSupport float32 // SLI Pro support
	CarPosition         float32 // car race position
	KersLevel           float32 // kers energy left
	KersMaxLevel        float32 // kers maximum energy
	Drs                 float32 // 0 = off, 1 = on
	TractionControl     float32 // 0 (off) - 2 (high)
	AntiLockBRakes      float32 // 0 (off) - 1 (on)
	FuelInTank          float32 // current fuel mass
	FuelCapacity        float32 // fuel capacity
	InPits              float32 // 0 = none, 1 = pitting, 2 = in pit area
	Sector              float32 // 0 = sector1, 1 = sector2; 2 = sector3
	Sector1Time         float32 // time of sector1 (or 0)
	Sector2Time         float32 // time of sector2 (or 0)
	BRakesTempBL        float32 // brakes temperature (centigrade)
	BRakesTempBR        float32 // brakes temperature (centigrade)
	BRakesTempFL        float32 // brakes temperature (centigrade)
	BRakesTempFR        float32 // brakes temperature (centigrade)
	TeamInfo            float32 // team ID
	TotalLaps           float32 // total number of laps in this race
	TrackSize           float32 // track size meters
	LastLapTime         float32 // last lap time
	MaxGears            float32 // maximum number of gears
	SessionType         float32 // 0 = unknown, 1 = practice, 2 = qualifying, 3 = race
	DRSAllowed          float32 // 0 = not allowed, 1 = allowed, -1 = invalid / unknown
	TrackNumber         float32 // -1 for unknown, 0-21 for tracks
	MaxRPM              float32 // cars max RPM, at which point the rev limiter will kick in
	IdleRPM             float32 // cars idle RPM
	VehicleFIAFlags     float32 // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
}

func ReadPacket(b []byte) (*TelemetryFrame, error) {
	buf := bytes.NewReader(b)

	frame := &TelemetryFrame{}
	err := binary.Read(buf, binary.LittleEndian, frame)
	if err != nil {
		return nil, err
	}

	return frame, nil
}

func TelemetryToDataFrame(tf TelemetryFrame) *data.Frame {
	frame := data.NewFrame("response")
	telemetryMap := telemetryFrameToMap(tf)

	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{time.Now()}),
	)

	for name, value := range telemetryMap {
		frame.Fields = append(frame.Fields,
			data.NewField(name, nil, []float32{value}),
		)
	}

	return frame
}

func telemetryFrameToMap(frame TelemetryFrame) map[string]float32 {
	var frameMap map[string]float32
	frame = convertTelemetryValues(frame)
	frameJson, err := json.Marshal(&frame)
	if err != nil {
		log.DefaultLogger.Error("Error converting frame", "error", err)
	}
	json.Unmarshal(frameJson, &frameMap)
	return frameMap
}

func convertTelemetryValues(f TelemetryFrame) TelemetryFrame {
	f.Speed = f.Speed * 3.6
	f.EngineRate = f.EngineRate * 10
	return f
}
