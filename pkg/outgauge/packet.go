package outgauge

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"time"
)

type OutgaugeStruct struct {
	Time          int32    // time in milliseconds (to check order)
	Car           [4]byte  // Car name
	Flags         [2]byte  // Info (see OG_x below)
	Gear          byte     // Reverse:0, Neutral:1, First:2...
	PlayerId      byte     // Unique ID of viewed player (0 = none)
	Speed         float32  // M/S
	RPM           float32  // RPM
	TurboPressure float32  // BAR
	EngTemp       float32  // C
	Fuel          float32  // 0 to 1
	OilPressure   float32  // BAR
	OilTemp       float32  // C
	DashLights    int32    // Dash lights available (see DL_x below)
	ShowLights    int32    // Dash lights currently switched on
	Throttle      float32  // 0 to 1
	Brake         float32  // 0 to 1
	Clutch        float32  // 0 to 1
	Display1      [16]byte // Usually Fuel
	Display2      [16]byte // Usually Settings
	Id            int32    // optional - only if OutGauge ID is specified
}

type DashLights struct {
	DL_SHIFT      int32
	DL_FULLBEAM   int32
	DL_HANDBRAKE  int32
	DL_PITSPEED   int32
	DL_TC         int32
	DL_SIGNAL_L   int32
	DL_SIGNAL_R   int32
	DL_SIGNAL_ANY int32
	DL_OILWARN    int32
	DL_BATTERY    int32
	DL_ABS        int32
	DL_SPARE      int32
}

// CONSTANTS
// // OG_x - bits for OutGaugePack Flags
// #define OG_SHIFT      1        // key
// #define OG_CTRL       2        // key
// #define OG_TURBO      8192     // show turbo gauge
// #define OG_KM         16384    // if not set - user prefers MILES
// #define OG_BAR        32768    // if not set - user prefers PSI
//
// // DL_x - bits for OutGaugePack DashLights and ShowLights
// DL_SHIFT,           // bit 0    - shift light
// DL_FULLBEAM,        // bit 1    - full beam
// DL_HANDBRAKE,       // bit 2    - handbrake
// DL_PITSPEED,        // bit 3    - pit speed limiter
// DL_TC,              // bit 4    - TC active or switched off
// DL_SIGNAL_L,        // bit 5    - left turn signal
// DL_SIGNAL_R,        // bit 6    - right turn signal
// DL_SIGNAL_ANY,      // bit 7    - shared turn signal
// DL_OILWARN,         // bit 8    - oil pressure warning
// DL_BATTERY,         // bit 9    - battery warning
// DL_ABS,             // bit 10   - ABS active or switched off
// DL_SPARE,           // bit 11

func ReadPacket(b []byte) (*OutgaugeStruct, error) {
	buf := bytes.NewReader(b)

	frame := &OutgaugeStruct{}
	err := binary.Read(buf, binary.LittleEndian, frame)
	if err != nil {
		return nil, err
	}

	return frame, nil
}

func TelemetryToDataFrame(t OutgaugeStruct) *data.Frame {
	t = convertTelemetryValues(t)
	frame := data.NewFrame("response")
	dl := readDashLights(t)

	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{time.Now()}),
		data.NewField("Car", nil, []string{string(t.Car[:])}),
		data.NewField("Display1", nil, []string{string(t.Display1[:])}),
		data.NewField("Display2", nil, []string{string(t.Display2[:])}),
		data.NewField("Gear", nil, []int8{int8(t.Gear)}),
		data.NewField("Speed", nil, []float32{t.Speed}),
		data.NewField("RPM", nil, []float32{t.RPM}),
		data.NewField("TurboPressure", nil, []float32{t.TurboPressure}),
		data.NewField("EngTemp", nil, []float32{t.EngTemp}),
		data.NewField("Fuel", nil, []float32{t.Fuel}),
		data.NewField("OilPressure", nil, []float32{t.OilPressure}),
		data.NewField("OilTemp", nil, []float32{t.OilTemp}),
		data.NewField("Throttle", nil, []float32{t.Throttle}),
		data.NewField("Brake", nil, []float32{t.Brake}),
		data.NewField("Clutch", nil, []float32{t.Clutch}),

		data.NewField("DL_SHIFT", nil, []int32{dl.DL_SHIFT}),
		data.NewField("DL_FULLBEAM", nil, []int32{dl.DL_FULLBEAM}),
		data.NewField("DL_HANDBRAKE", nil, []int32{dl.DL_HANDBRAKE}),
		data.NewField("DL_PITSPEED", nil, []int32{dl.DL_PITSPEED}),
		data.NewField("DL_TC", nil, []int32{dl.DL_TC}),
		data.NewField("DL_SIGNAL_L", nil, []int32{dl.DL_SIGNAL_L}),
		data.NewField("DL_SIGNAL_R", nil, []int32{dl.DL_SIGNAL_R}),
		data.NewField("DL_SIGNAL_ANY", nil, []int32{dl.DL_SIGNAL_ANY}),
		data.NewField("DL_OILWARN", nil, []int32{dl.DL_OILWARN}),
		data.NewField("DL_BATTERY", nil, []int32{dl.DL_BATTERY}),
		data.NewField("DL_ABS", nil, []int32{dl.DL_ABS}),
		data.NewField("DL_SPARE", nil, []int32{dl.DL_SPARE}),
	)

	return frame
}

func telemetryFrameToMap(frame OutgaugeStruct) map[string]float32 {
	var frameMap map[string]float32
	frame = convertTelemetryValues(frame)
	frameJson, err := json.Marshal(&frame)
	if err != nil {
		log.DefaultLogger.Error("Error converting frame", "error", err)
	}
	json.Unmarshal(frameJson, &frameMap)
	return frameMap
}

func convertTelemetryValues(f OutgaugeStruct) OutgaugeStruct {
	f.Speed = f.Speed * 3.6
	f.Gear = f.Gear - 1
	return f
}

func readDashLights(f OutgaugeStruct) DashLights {
	dl := DashLights{}
	dlData := int32(f.ShowLights)
	if dlData&1 != 0 {
		dl.DL_SHIFT = 1
	}
	if dlData&2 != 0 {
		dl.DL_FULLBEAM = 1
	}
	if dlData&4 != 0 {
		dl.DL_HANDBRAKE = 1
	}
	if dlData&8 != 0 {
		dl.DL_PITSPEED = 1
	}
	if dlData&16 != 0 {
		dl.DL_TC = 1
	}
	if dlData&32 != 0 {
		dl.DL_SIGNAL_L = 1
	}
	if dlData&64 != 0 {
		dl.DL_SIGNAL_R = 1
	}
	if dlData&128 != 0 {
		dl.DL_SIGNAL_ANY = 1
	}
	if dlData&256 != 0 {
		dl.DL_OILWARN = 1
	}
	if dlData&512 != 0 {
		dl.DL_BATTERY = 1
	}
	if dlData&1024 != 0 {
		dl.DL_ABS = 1
	}
	if dlData&2048 != 0 {
		dl.DL_SPARE = 1
	}
	return dl
}
