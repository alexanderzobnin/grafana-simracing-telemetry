package sharedmemory

import (
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"time"
)

type SPageFilePhysics struct {
	Packetid            int32
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

type TelemetryPhysics struct {
	Packetid            int32
	Gas                 float32
	Brake               float32
	Fuel                float32
	Gear                int32
	Rpms                int32
	Steerangle          float32
	Speedkmh            float32
	Velocity            [3]float32
	Accg                [3]float32
	Wheelslip           [4]float32
	Wheelload           [4]float32
	Wheelspressure      [4]float32
	Wheelangularspeed   [4]float32
	Tyrewear            [4]float32
	Tyredirtylevel      [4]float32
	Tyrecoretemperature [4]float32
	Camberrad           [4]float32
	Suspensiontravel    [4]float32
	Drs                 float32
	Tc                  float32
	Heading             float32
	Pitch               float32
	Roll                float32
	Cgheight            float32
	Cardamage           [5]float32
	Numberoftyresout    int32
	Pitlimiteron        int32
	Abs                 float32
	Kerscharge          float32
	Kersinput           float32
	Autoshifteron       int32
	Rideheight          [2]float32
	Turboboost          float32
	Ballast             float32
	Airdensity          float32
	Airtemp             float32
	Roadtemp            float32
	Localangularvel     [3]float32
	Finalff             float32
	Performancemeter    float32
	Enginebrake         int32
	Ersrecoverylevel    int32
	Erspowerlevel       int32
	Ersheatcharging     int32
	Ersischarging       int32
	Kerscurrentkj       float32
	Drsavailable        int32
	Drsenabled          int32
	Braketemp           [4]float32
	Clutch              float32
	Tyretempi           [4]float32
	Tyretempm           [4]float32
	Tyretempo           [4]float32
	Isaicontrolled      int32
	Tyrecontactpoint32  [4][3]float32
	Tyrecontactnormal   [4][3]float32
	Tyrecontactheading  [4][3]float32
	Brakebias           float32
	Localvelocity       [3]float32
	P2pactivations      int32
	P2pstatus           int32
	Currentmaxrpm       int32
	Mz                  [4]float32
	Fx                  [4]float32
	Fy                  [4]float32
	Slipratio           [4]float32
	Slipangle           [4]float32
	Tcinaction          int32
	Absinaction         int32
	Suspensiondamage    [4]float32
	Tyretemp            [4]float32
}

func (f *SPageFilePhysics) convertTelemetryValues() *SPageFilePhysics {
	f.Gear = f.Gear - 1
	return f
}

func PhysicsToDataFrame(t SPageFilePhysics) (*data.Frame, error) {
	t.convertTelemetryValues()

	frame := data.NewFrame("response")
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{time.Now()}),
	)

	frame.Fields = append(frame.Fields,
		data.NewField("Packetid", nil, []int32{t.Packetid}),
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
