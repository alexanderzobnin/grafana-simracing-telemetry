package sharedmemory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"golang.org/x/text/encoding/unicode"
	"time"
)

// SPageFilePhysics updates at each physics step. They all refer to the player’s car.
type SPageFilePhysics struct {
	PacketId            int32
	Gas                 float32
	Brake               float32
	Fuel                float32
	Gear                int32
	RPMs                int32
	SteerAngle          float32
	SpeedKmh            float32
	Velocity            [3]float32
	AccG                [3]float32
	WheelSlip           [4]float32
	WheelLoad           [4]float32 // Field is not used by ACC
	WheelsPressure      [4]float32
	WheelAngularSpeed   [4]float32
	TyreWear            [4]float32 // Field is not used by ACC
	TyreDirtyLevel      [4]float32 // Field is not used by ACC
	TyreCoreTemperature [4]float32
	CamberRad           [4]float32 // Field is not used by ACC
	SuspensionTravel    [4]float32
	DRS                 float32 // Field is not used by ACC
	TC                  float32
	Heading             float32
	Pitch               float32
	Roll                float32
	CGHeight            float32    // Field is not used by ACC
	CarDamage           [5]float32 // Car damage: front 0, rear 1, left 2, right 3, centre 4
	NumberOfTyresOut    int32      // Field is not used by ACC
	PitLimiterOn        int32
	ABS                 float32
	KersCharge          float32 // Field is not used by ACC
	KersInput           float32 // Field is not used by ACC
	AutoShifterOn       int32
	RideHeight          [2]float32 // Field is not used by ACC
	TurboBoost          float32
	Ballast             float32 // Field is not used by ACC
	AirDensity          float32 // Field is not used by ACC
	AirTemp             float32
	RoadTemp            float32
	LocalAngularVel     [3]float32
	FinalFF             float32 // Force feedback signal
	PerformanceMeter    float32 // Field is not used by ACC
	Enginebrake         int32   // Field is not used by ACC
	Ersrecoverylevel    int32   // Field is not used by ACC
	Erspowerlevel       int32   // Field is not used by ACC
	Ersheatcharging     int32   // Field is not used by ACC
	Ersischarging       int32   // Field is not used by ACC
	Kerscurrentkj       float32 // Field is not used by ACC
	Drsavailable        int32   // Field is not used by ACC
	Drsenabled          int32   // Field is not used by ACC
	BrakeTemp           [4]float32
	Clutch              float32
	TyretempI           [4]float32 // Field is not used by ACC
	TyretempM           [4]float32 // Field is not used by ACC
	TyretempO           [4]float32 // Field is not used by ACC
	IsAIControlled      int32
	TyreContactPoint    [4][3]float32 // Tyre contact point global coordinates [FL, FR, RL, RR
	TyreContactNormal   [4][3]float32 // Tyre contact normal [FL, FR, RL, RR] [x,y,z]
	TyreContactHeading  [4][3]float32 // Tyre contact heading [FL, FR, RL, RR] [x,y,z]
	BrakeBias           float32       // Front brake bias, see Appendix 4
	LocalVelocity       [3]float32
	P2pactivations      int32      // Field is not used by ACC
	P2pstatus           int32      // Field is not used by ACC
	Currentmaxrpm       int32      // Field is not used by ACC
	Mz                  [4]float32 // Field is not used by ACC
	Fx                  [4]float32 // Field is not used by ACC
	Fy                  [4]float32 // Field is not used by ACC
	SlipRatio           [4]float32 // Tyre slip ratio [FL, FR, RL, RR] in radians
	SlipAngle           [4]float32 // Tyre slip angle [FL, FR, RL, RR]
	TCInAction          int32      // Field is not used by ACC
	ABSInAction         int32      // Field is not used by ACC
	SuspensionDamage    [4]float32 // Field is not used by ACC
	TyreTemp            [4]float32 // Field is not used by ACC
}

type SPageFileGraphic struct {
	PacketId                 int32
	Status                   int32
	Session                  int32
	CurrentTime              [30]byte
	LastTime                 [30]byte
	BestTime                 [30]byte
	Split                    [30]byte
	CompletedLaps            int32
	Position                 int32
	ICurrentTime             int32
	ILastTime                int32
	IBestTime                int32
	SessionTimeLeft          float32
	DistanceTraveled         float32
	IsInPit                  int32
	CurrentSectorIndex       int32
	LastSectorTime           int32
	NumberOfLaps             int32
	TyreCompound             [68]byte
	ReplayTimeMultiplier     float32
	NormalizedCarPosition    float32
	ActiveCars               int32
	CarCoordinates           [60][3]float32
	CarId                    [60]int32
	PlayercarId              int32
	PenaltyTime              float32
	Flag                     int32
	Penalty                  int32
	IdealLineOn              int32
	IsInPitLane              int32
	SurfaceGrip              float32
	MandatoryPitDone         int32
	WindSpeed                float32
	WindDirection            float32
	IsSetupMenuVisible       int32
	MainDisplayIndex         int32
	SecondaryDisplayIndex    int32
	TCLevel                  int32
	TCCut                    int32
	EngineMap                int32
	ABSLevel                 int32
	FuelxLap                 float32
	RainLights               int32
	FlashingLights           int32
	LightsStage              int32
	ExhaustTemperature       float32
	WiperLv                  int32
	DriverStintTotalTimeLeft int32
	DriverStintTimeLeft      int32
	RainTyres                int32

	SessionIndex         int32
	UsedFuel             float32
	DeltaLapTime         [32]byte
	IDeltaLapTime        int32
	EstimatedLapTime     [32]byte
	IEstimatedLapTime    int32
	IsDeltaPositive      int32
	ISplit               int32
	IsValidLap           int32
	FuelEstimatedLaps    float32
	TrackStatus          [68]byte
	MissingMandatoryPits int32
	Clock                float32
}

type ACCTelemetry struct {
	SPageFileGraphic
	SPageFilePhysics
}

func (f *SPageFilePhysics) convertTelemetryValues() *SPageFilePhysics {
	f.Gear = f.Gear - 1
	return f
}

func ACCTelemetryToDataFrame(t ACCTelemetry) (*data.Frame, error) {
	frame, err := PhysicsToDataFrame(t.SPageFilePhysics)
	if err != nil {
		return nil, err
	}
	return GraphicToDataFrame(t.SPageFileGraphic, frame)
}

func PhysicsToDataFrame(t SPageFilePhysics) (*data.Frame, error) {
	t.convertTelemetryValues()

	frame := data.NewFrame("response")
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{time.Now()}),
	)

	frame.Fields = append(frame.Fields,
		data.NewField("PacketId", nil, []int32{t.PacketId}),
		data.NewField("Gas", nil, []float32{t.Gas}),
		data.NewField("Brake", nil, []float32{t.Brake}),
		data.NewField("Fuel", nil, []float32{t.Fuel}),
		data.NewField("Gear", nil, []int32{t.Gear}),
		data.NewField("RPMs", nil, []int32{t.RPMs}),
		data.NewField("SteerAngle", nil, []float32{t.SteerAngle}),
		data.NewField("SpeedKmh", nil, []float32{t.SpeedKmh}),
		data.NewField("TC", nil, []float32{t.TC}),
		data.NewField("Heading", nil, []float32{t.Heading}),
		data.NewField("Pitch", nil, []float32{t.Pitch}),
		data.NewField("Roll", nil, []float32{t.Roll}),
		data.NewField("ABS", nil, []float32{t.ABS}),
		data.NewField("TurboBoost", nil, []float32{t.TurboBoost}),
		data.NewField("AirTemp", nil, []float32{t.AirTemp}),
		data.NewField("RoadTemp", nil, []float32{t.RoadTemp}),
		data.NewField("FinalFF", nil, []float32{t.FinalFF}),
		data.NewField("Clutch", nil, []float32{t.Clutch}),
		data.NewField("BrakeBias", nil, []float32{t.BrakeBias}),
		data.NewField("AutoShifterOn", nil, []int32{t.AutoShifterOn}),
		data.NewField("PitLimiterOn", nil, []int32{t.PitLimiterOn}),
		data.NewField("IsAIControlled", nil, []int32{t.IsAIControlled}),
	)

	frame.Fields = append(frame.Fields, toTiresFields("BrakeTemp", t.BrakeTemp)...)
	frame.Fields = append(frame.Fields, toTiresFields("SlipAngle", t.SlipAngle)...)
	frame.Fields = append(frame.Fields, toTiresFields("SlipRatio", t.SlipRatio)...)
	frame.Fields = append(frame.Fields, toTiresFields("WheelSlip", t.WheelSlip)...)
	frame.Fields = append(frame.Fields, toTiresFields("TyreCoreTemperature", t.TyreCoreTemperature)...)
	frame.Fields = append(frame.Fields, toTiresFields("WheelsPressure", t.WheelsPressure)...)
	frame.Fields = append(frame.Fields, toTiresFields("SuspensionTravel", t.SuspensionTravel)...)
	frame.Fields = append(frame.Fields, toCoordinatesFields("LocalAngularVel", t.LocalAngularVel)...)
	frame.Fields = append(frame.Fields, toCoordinatesFields("LocalVelocity", t.LocalVelocity)...)
	frame.Fields = append(frame.Fields, toCoordinatesFields("Velocity", t.Velocity)...)
	frame.Fields = append(frame.Fields, toCoordinatesFields("AccG", t.AccG)...)
	frame.Fields = append(frame.Fields, toTiresCoordinatesFields("TyreContactPoint", t.TyreContactPoint)...)
	frame.Fields = append(frame.Fields, toTiresCoordinatesFields("TyreContactNormal", t.TyreContactNormal)...)
	frame.Fields = append(frame.Fields, toTiresCoordinatesFields("TyreContactHeading", t.TyreContactHeading)...)

	return frame, nil
}

func GraphicToDataFrame(t SPageFileGraphic, frame *data.Frame) (*data.Frame, error) {
	frame.Fields = append(frame.Fields,
		data.NewField("Status", nil, []int32{t.Status}),
		data.NewField("Session", nil, []int32{t.Session}),
		data.NewField("CurrentTime", nil, []string{wchart30ToString(t.CurrentTime)}),
		data.NewField("LastTime", nil, []string{wchart30ToString(t.LastTime)}),
		data.NewField("BestTime", nil, []string{wchart30ToString(t.BestTime)}),
		data.NewField("Split", nil, []string{wchart30ToString(t.Split)}),
		data.NewField("CompletedLaps", nil, []int32{t.CompletedLaps}),
		data.NewField("Position", nil, []int32{t.Position}),
		data.NewField("ICurrentTime", nil, []int32{t.ICurrentTime}),
		data.NewField("ILastTime", nil, []int32{t.ILastTime}),
		data.NewField("IBestTime", nil, []int32{t.IBestTime}),
		data.NewField("SessionTimeLeft", nil, []float32{t.SessionTimeLeft}),
		data.NewField("DistanceTraveled", nil, []float32{t.DistanceTraveled}),
		data.NewField("IBestTime", nil, []int32{t.IsInPit}),
		data.NewField("CurrentSectorIndex", nil, []int32{t.CurrentSectorIndex}),
		data.NewField("LastSectorTime", nil, []int32{t.LastSectorTime}),
		data.NewField("NumberOfLaps", nil, []int32{t.NumberOfLaps}),
		data.NewField("TyreCompound", nil, []string{wchart68ToString(t.TyreCompound)}),
		data.NewField("NormalizedCarPosition", nil, []float32{t.NormalizedCarPosition}),
		//data.NewField("CarCoordinates", nil, []float32{t.CarCoordinates}),
		data.NewField("ActiveCars", nil, []int32{t.ActiveCars}),
		//data.NewField("CarId", nil, []int32{t.CarId}),
		data.NewField("PlayercarId", nil, []int32{t.PlayercarId}),
		data.NewField("PenaltyTime", nil, []float32{t.PenaltyTime}),
		data.NewField("Flag", nil, []int32{t.Flag}),
		data.NewField("Penalty", nil, []int32{t.Penalty}),
		data.NewField("IdealLineOn", nil, []int32{t.IdealLineOn}),
		data.NewField("IsInPitLane", nil, []int32{t.IsInPitLane}),
		data.NewField("SurfaceGrip", nil, []float32{t.SurfaceGrip}),
		data.NewField("MandatoryPitDone", nil, []int32{t.MandatoryPitDone}),
		data.NewField("WindSpeed", nil, []float32{t.WindSpeed}),
		data.NewField("WindDirection", nil, []float32{t.WindDirection}),
		data.NewField("IsSetupMenuVisible", nil, []int32{t.IsSetupMenuVisible}),
		data.NewField("MainDisplayIndex", nil, []int32{t.MainDisplayIndex}),
		data.NewField("SecondaryDisplayIndex", nil, []int32{t.SecondaryDisplayIndex}),
		data.NewField("TCLevel", nil, []int32{t.TCLevel}),
		data.NewField("TCCut", nil, []int32{t.TCCut}),
		data.NewField("EngineMap", nil, []int32{t.EngineMap}),
		data.NewField("ABSLevel", nil, []int32{t.ABSLevel}),
		data.NewField("FuelxLap", nil, []float32{t.FuelxLap}),
		data.NewField("RainLights", nil, []int32{t.RainLights}),
		data.NewField("FlashingLights", nil, []int32{t.FlashingLights}),
		data.NewField("LightsStage", nil, []int32{t.LightsStage}),
		data.NewField("ExhaustTemperature", nil, []float32{t.ExhaustTemperature}),
		data.NewField("WiperLv", nil, []int32{t.WiperLv}),
		data.NewField("DriverStintTotalTimeLeft", nil, []int32{t.DriverStintTotalTimeLeft}),
		data.NewField("DriverStintTimeLeft", nil, []int32{t.DriverStintTimeLeft}),
		data.NewField("RainTyres", nil, []int32{t.RainTyres}),

		data.NewField("SessionIndex", nil, []int32{t.SessionIndex}),
		data.NewField("UsedFuel", nil, []float32{t.UsedFuel}),
		data.NewField("DeltaLapTime", nil, []string{wchart32ToString(t.DeltaLapTime)}),
		data.NewField("IDeltaLapTime", nil, []int32{t.IDeltaLapTime}),
		data.NewField("EstimatedLapTime", nil, []string{wchart32ToString(t.EstimatedLapTime)}),
		data.NewField("IEstimatedLapTime", nil, []int32{t.IEstimatedLapTime}),
		data.NewField("IsDeltaPositive", nil, []int32{t.IsDeltaPositive}),
		data.NewField("ISplit", nil, []int32{t.ISplit}),
		data.NewField("IsValidLap", nil, []int32{t.IsValidLap}),
		data.NewField("FuelEstimatedLaps", nil, []float32{t.FuelEstimatedLaps}),
		data.NewField("TrackStatus", nil, []string{wchart68ToString(t.TrackStatus)}),
		data.NewField("MissingMandatoryPits", nil, []int32{t.MissingMandatoryPits}),
		data.NewField("Clock", nil, []float32{t.Clock}),
	)

	return frame, nil
}

func toTiresFields(name string, values [4]float32) []*data.Field {
	return []*data.Field{
		data.NewField(fmt.Sprintf("%sFL", name), nil, []float64{float64(values[0])}),
		data.NewField(fmt.Sprintf("%sFR", name), nil, []float64{float64(values[1])}),
		data.NewField(fmt.Sprintf("%sRL", name), nil, []float64{float64(values[2])}),
		data.NewField(fmt.Sprintf("%sRR", name), nil, []float64{float64(values[3])}),
	}
}

func toCoordinatesFields(name string, values [3]float32) []*data.Field {
	return []*data.Field{
		data.NewField(fmt.Sprintf("%sX", name), nil, []float64{float64(values[0])}),
		data.NewField(fmt.Sprintf("%sY", name), nil, []float64{float64(values[1])}),
		data.NewField(fmt.Sprintf("%sZ", name), nil, []float64{float64(values[2])}),
	}
}

func toTiresCoordinatesFields(name string, values [4][3]float32) []*data.Field {
	return []*data.Field{
		data.NewField(fmt.Sprintf("%sFLX", name), nil, []float64{float64(values[0][0])}),
		data.NewField(fmt.Sprintf("%sFLY", name), nil, []float64{float64(values[0][1])}),
		data.NewField(fmt.Sprintf("%sFLZ", name), nil, []float64{float64(values[0][2])}),

		data.NewField(fmt.Sprintf("%sFRX", name), nil, []float64{float64(values[1][0])}),
		data.NewField(fmt.Sprintf("%sFRY", name), nil, []float64{float64(values[1][1])}),
		data.NewField(fmt.Sprintf("%sFRZ", name), nil, []float64{float64(values[1][2])}),

		data.NewField(fmt.Sprintf("%sRLX", name), nil, []float64{float64(values[2][0])}),
		data.NewField(fmt.Sprintf("%sRLY", name), nil, []float64{float64(values[2][1])}),
		data.NewField(fmt.Sprintf("%sRLZ", name), nil, []float64{float64(values[2][2])}),

		data.NewField(fmt.Sprintf("%sRRX", name), nil, []float64{float64(values[3][0])}),
		data.NewField(fmt.Sprintf("%sRRY", name), nil, []float64{float64(values[3][1])}),
		data.NewField(fmt.Sprintf("%sRRZ", name), nil, []float64{float64(values[3][2])}),
	}
}

func telemetryToMap(frame SPageFilePhysics) (map[string]interface{}, error) {
	var frameMap map[string]interface{}
	frameJson, err := json.Marshal(&frame)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(frameJson, &frameMap)
	return frameMap, nil
}

func wchartToString(b []byte) string {
	dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	out, err := dec.Bytes(b)
	if err != nil {
		return ""
	}

	// Strings are null terminated
	i := bytes.IndexByte(out, 0)
	if i == -1 {
		i = len(out)
	}

	return string(out[:i])
}

func wchart30ToString(bFixed [30]byte) string {
	b := make([]byte, 0)
	for _, v := range bFixed {
		b = append(b, v)
	}
	return wchartToString(b)
}

func wchart68ToString(bFixed [68]byte) string {
	b := make([]byte, 0)
	for _, v := range bFixed {
		b = append(b, v)
	}
	return wchartToString(b)
}

func wchart32ToString(bFixed [32]byte) string {
	b := make([]byte, 0)
	for _, v := range bFixed {
		b = append(b, v)
	}
	return wchartToString(b)
}
