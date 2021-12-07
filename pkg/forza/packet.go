package forza

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type TelemetryFrame struct {
	IsRaceOn                             int32  // = 1 when race is on. = 0 when in menus/race stopped … // s32
	TimestampMS                          uint32 //Can overflow to 0 eventually
	EngineMaxRpm                         float32
	EngineIdleRpm                        float32
	CurrentEngineRpm                     float32
	AccelerationX                        float32 //In the car's local space  X = right, Y = up, Z = forward
	AccelerationY                        float32
	AccelerationZ                        float32
	VelocityX                            float32 //In the car's local space  X = right, Y = up, Z = forward
	VelocityY                            float32
	VelocityZ                            float32
	AngularVelocityX                     float32 //In the car's local space  X = pitch, Y = yaw, Z = roll
	AngularVelocityY                     float32
	AngularVelocityZ                     float32
	Yaw                                  float32
	Pitch                                float32
	Roll                                 float32
	NormalizedSuspensionTravelFrontLeft  float32 // Suspension travel normalized: 0.0f = max stretch  1.0 = max compression
	NormalizedSuspensionTravelFrontRight float32
	NormalizedSuspensionTravelRearLeft   float32
	NormalizedSuspensionTravelRearRight  float32
	TireSlipRatioFrontLeft               float32 // Tire normalized slip ratio, = 0 means 100% grip and |ratio| > 1.0 means loss of grip.
	TireSlipRatioFrontRight              float32
	TireSlipRatioRearLeft                float32
	TireSlipRatioRearRight               float32
	WheelRotationSpeedFrontLeft          float32 // Wheel rotation speed radians/sec.
	WheelRotationSpeedFrontRight         float32
	WheelRotationSpeedRearLeft           float32
	WheelRotationSpeedRearRight          float32
	WheelOnRumbleStripFrontLeft          uint32  // = 1 when wheel is on rumble strip, = 0 when off. // s32
	WheelOnRumbleStripFrontRight         uint32  // s32
	WheelOnRumbleStripRearLeft           uint32  // s32
	WheelOnRumbleStripRearRight          uint32  // s32
	WheelInPuddleDepthFrontLeft          float32 // = from 0 to 1, where 1 is the deepest puddle
	WheelInPuddleDepthFrontRight         float32
	WheelInPuddleDepthRearLeft           float32
	WheelInPuddleDepthRearRight          float32
	SurfaceRumbleFrontLeft               float32 // Non-dimensional surface rumble values passed to controller force feedback
	SurfaceRumbleFrontRight              float32
	SurfaceRumbleRearLeft                float32
	SurfaceRumbleRearRight               float32
	TireSlipAngleFrontLeft               float32 // Tire normalized slip angle, = 0 means 100% grip and |angle| > 1.0 means loss of grip.
	TireSlipAngleFrontRight              float32
	TireSlipAngleRearLeft                float32
	TireSlipAngleRearRight               float32
	TireCombinedSlipFrontLeft            float32 // Tire normalized combined slip, = 0 means 100% grip and |slip| > 1.0 means loss of grip.
	TireCombinedSlipFrontRight           float32
	TireCombinedSlipRearLeft             float32
	TireCombinedSlipRearRight            float32
	SuspensionTravelMetersFrontLeft      float32 // Actual suspension travel in meters
	SuspensionTravelMetersFrontRight     float32
	SuspensionTravelMetersRearLeft       float32
	SuspensionTravelMetersRearRight      float32
	CarOrdinal                           int32 //Unique ID of the car make/model // s32
	CarClass                             int32 //Between 0 (D -- worst cars) and 7 (X class -- best cars) inclusive // s32
	CarPerformanceIndex                  int32 //Between 100 (slowest car) and 999 (fastest car) inclusive // s32
	DrivetrainType                       int32 //Corresponds to EDrivetrainType; 0 = FWD, 1 = RWD, 2 = AWD // s32
	NumCylinders                         int32 //Number of cylinders in the engine // s32
	CarCategory                          int32
	HorizonPlaceholder1                  int32   // > 0 when crashing into objects
	HorizonPlaceholder2                  int32   // > 0 when crashing into objects
	PositionX                            float32 // in meters
	PositionY                            float32
	PositionZ                            float32
	Speed                                float32 // meters per second
	Power                                float32 // watts
	Torque                               float32 // newton meter
	TireTempFrontLeft                    float32 // °F
	TireTempFrontRight                   float32
	TireTempRearLeft                     float32
	TireTempRearRight                    float32
	Boost                                float32
	Fuel                                 float32
	DistanceTraveled                     float32
	BestLap                              float32
	LastLap                              float32
	CurrentLap                           float32
	CurrentRaceTime                      float32
	LapNumber                            uint16
	RacePosition                         uint8
	Throttle                             uint8
	Brake                                uint8
	Clutch                               uint8
	HandBrake                            uint8
	Gear                                 uint8
	Steer                                uint8
	NormalizedDrivingLine                int8
	NormalizedAIBrakeDifference          int8
}

type TelemetryFrameConverted struct {
	IsRaceOn                             int32  // = 1 when race is on. = 0 when in menus/race stopped … // s32
	TimestampMS                          uint32 //Can overflow to 0 eventually
	EngineMaxRpm                         float32
	EngineIdleRpm                        float32
	CurrentEngineRpm                     float32
	AccelerationX                        float32 //In the car's local space  X = right, Y = up, Z = forward
	AccelerationY                        float32
	AccelerationZ                        float32
	VelocityX                            float32 //In the car's local space  X = right, Y = up, Z = forward
	VelocityY                            float32
	VelocityZ                            float32
	AngularVelocityX                     float32 //In the car's local space  X = pitch, Y = yaw, Z = roll
	AngularVelocityY                     float32
	AngularVelocityZ                     float32
	Yaw                                  float32
	Pitch                                float32
	Roll                                 float32
	NormalizedSuspensionTravelFrontLeft  float32 // Suspension travel normalized: 0.0f = max stretch  1.0 = max compression
	NormalizedSuspensionTravelFrontRight float32
	NormalizedSuspensionTravelRearLeft   float32
	NormalizedSuspensionTravelRearRight  float32
	TireSlipRatioFrontLeft               float32 // Tire normalized slip ratio, = 0 means 100% grip and |ratio| > 1.0 means loss of grip.
	TireSlipRatioFrontRight              float32
	TireSlipRatioRearLeft                float32
	TireSlipRatioRearRight               float32
	WheelRotationSpeedFrontLeft          float32 // Wheel rotation speed radians/sec.
	WheelRotationSpeedFrontRight         float32
	WheelRotationSpeedRearLeft           float32
	WheelRotationSpeedRearRight          float32
	WheelOnRumbleStripFrontLeft          uint32  // = 1 when wheel is on rumble strip, = 0 when off. // s32
	WheelOnRumbleStripFrontRight         uint32  // s32
	WheelOnRumbleStripRearLeft           uint32  // s32
	WheelOnRumbleStripRearRight          uint32  // s32
	WheelInPuddleDepthFrontLeft          float32 // = from 0 to 1, where 1 is the deepest puddle
	WheelInPuddleDepthFrontRight         float32
	WheelInPuddleDepthRearLeft           float32
	WheelInPuddleDepthRearRight          float32
	SurfaceRumbleFrontLeft               float32 // Non-dimensional surface rumble values passed to controller force feedback
	SurfaceRumbleFrontRight              float32
	SurfaceRumbleRearLeft                float32
	SurfaceRumbleRearRight               float32
	TireSlipAngleFrontLeft               float32 // Tire normalized slip angle, = 0 means 100% grip and |angle| > 1.0 means loss of grip.
	TireSlipAngleFrontRight              float32
	TireSlipAngleRearLeft                float32
	TireSlipAngleRearRight               float32
	TireCombinedSlipFrontLeft            float32 // Tire normalized combined slip, = 0 means 100% grip and |slip| > 1.0 means loss of grip.
	TireCombinedSlipFrontRight           float32
	TireCombinedSlipRearLeft             float32
	TireCombinedSlipRearRight            float32
	SuspensionTravelMetersFrontLeft      float32 // Actual suspension travel in meters
	SuspensionTravelMetersFrontRight     float32
	SuspensionTravelMetersRearLeft       float32
	SuspensionTravelMetersRearRight      float32
	CarOrdinal                           int32 //Unique ID of the car make/model // s32
	CarClass                             int32 //Between 0 (D -- worst cars) and 7 (X class -- best cars) inclusive // s32
	CarPerformanceIndex                  int32 //Between 100 (slowest car) and 999 (fastest car) inclusive // s32
	DrivetrainType                       int32 //Corresponds to EDrivetrainType; 0 = FWD, 1 = RWD, 2 = AWD // s32
	NumCylinders                         int32 //Number of cylinders in the engine // s32
	CarCategory                          int32
	HorizonPlaceholder1                  int32   // > 0 when crashing into objects
	HorizonPlaceholder2                  int32   // > 0 when crashing into objects
	PositionX                            float32 // in meters
	PositionY                            float32
	PositionZ                            float32
	Speed                                float32 // meters per second
	Power                                float32 // watts
	Torque                               float32 // newton meter
	TireTempFrontLeft                    float32 // °F
	TireTempFrontRight                   float32
	TireTempRearLeft                     float32
	TireTempRearRight                    float32
	Boost                                float32
	Fuel                                 float32
	DistanceTraveled                     float32
	BestLap                              float32
	LastLap                              float32
	CurrentLap                           float32
	CurrentRaceTime                      float32
	LapNumber                            uint16
	RacePosition                         uint8
	Throttle                             float32
	Brake                                float32
	Clutch                               float32
	HandBrake                            float32
	Gear                                 uint8
	Steer                                float32
	NormalizedDrivingLine                int8
	NormalizedAIBrakeDifference          int8
	CarAttitude                          int8
	IsTractionLost                       int8
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
	frameConverted := convertTelemetryValues(frame)
	frameJson, err := json.Marshal(&frameConverted)
	if err != nil {
		log.DefaultLogger.Error("Error converting frame", "error", err)
	}
	json.Unmarshal(frameJson, &frameMap)
	return frameMap
}

func toCelsius(F float32) float32 {
	if F == 0 {
		return F
	}
	return (F - 32) / 1.8
}

// CheckAttitude looks for balance of the car
func CheckAttitude(totalSlipFront int8, totalSlipRear int8) (int8, int8) {
	carAttitude := int8(0)
	isTractionLost := int8(0)

	// Check attitude of car by comparing front and rear slip levels
	// If front slip > rear slip, means car is understeering
	if totalSlipRear > totalSlipFront {
		carAttitude = 2 // "Oversteer"
	} else if totalSlipFront > totalSlipRear {
		carAttitude = 1 // "Understeer"
	}

	if (totalSlipRear+totalSlipFront) > 2 && carAttitude == 2 {
		isTractionLost = 1
	}

	return carAttitude, isTractionLost

}

func convertTelemetryValues(f TelemetryFrame) TelemetryFrameConverted {
	mps_to_kmh := float32(3.6)
	watts_to_bhp := float32(1 / 745.7)
	totalSlipRear := int8(f.TireCombinedSlipRearLeft + f.TireCombinedSlipRearRight)
	totalSlipFront := int8(f.TireCombinedSlipFrontLeft + f.TireCombinedSlipFrontRight)
	carAttitude, isTractionLost := CheckAttitude(totalSlipFront, totalSlipRear)

	fC := TelemetryFrameConverted{
		IsRaceOn:                             f.IsRaceOn,
		TimestampMS:                          f.TimestampMS,
		EngineMaxRpm:                         f.EngineMaxRpm,
		EngineIdleRpm:                        f.EngineIdleRpm,
		CurrentEngineRpm:                     f.CurrentEngineRpm,
		AccelerationX:                        f.AccelerationX,
		AccelerationY:                        f.AccelerationY,
		AccelerationZ:                        f.AccelerationZ,
		VelocityX:                            f.VelocityX,
		VelocityY:                            f.VelocityY,
		VelocityZ:                            f.VelocityZ,
		AngularVelocityX:                     f.AngularVelocityX,
		AngularVelocityY:                     f.AngularVelocityY,
		AngularVelocityZ:                     f.AngularVelocityZ,
		Yaw:                                  f.Yaw,
		Pitch:                                f.Pitch,
		Roll:                                 f.Roll,
		NormalizedSuspensionTravelFrontLeft:  f.NormalizedSuspensionTravelFrontLeft,
		NormalizedSuspensionTravelFrontRight: f.NormalizedSuspensionTravelFrontRight,
		NormalizedSuspensionTravelRearLeft:   f.NormalizedSuspensionTravelRearLeft,
		NormalizedSuspensionTravelRearRight:  f.NormalizedSuspensionTravelRearRight,
		TireSlipRatioFrontLeft:               f.TireSlipRatioFrontLeft,
		TireSlipRatioFrontRight:              f.TireSlipRatioFrontRight,
		TireSlipRatioRearLeft:                f.TireSlipRatioRearLeft,
		TireSlipRatioRearRight:               f.TireSlipRatioRearRight,
		WheelRotationSpeedFrontLeft:          f.WheelRotationSpeedFrontLeft,
		WheelRotationSpeedFrontRight:         f.WheelRotationSpeedFrontRight,
		WheelRotationSpeedRearLeft:           f.WheelRotationSpeedRearLeft,
		WheelRotationSpeedRearRight:          f.WheelRotationSpeedRearRight,
		WheelOnRumbleStripFrontLeft:          f.WheelOnRumbleStripFrontLeft,
		WheelOnRumbleStripFrontRight:         f.WheelOnRumbleStripFrontRight,
		WheelOnRumbleStripRearLeft:           f.WheelOnRumbleStripRearLeft,
		WheelOnRumbleStripRearRight:          f.WheelOnRumbleStripRearRight,
		WheelInPuddleDepthFrontLeft:          f.WheelInPuddleDepthFrontLeft,
		WheelInPuddleDepthFrontRight:         f.WheelInPuddleDepthFrontRight,
		WheelInPuddleDepthRearLeft:           f.WheelInPuddleDepthRearLeft,
		WheelInPuddleDepthRearRight:          f.WheelInPuddleDepthRearRight,
		SurfaceRumbleFrontLeft:               f.SurfaceRumbleFrontLeft,
		SurfaceRumbleFrontRight:              f.SurfaceRumbleFrontRight,
		SurfaceRumbleRearLeft:                f.SurfaceRumbleRearLeft,
		SurfaceRumbleRearRight:               f.SurfaceRumbleRearRight,
		TireSlipAngleFrontLeft:               f.TireSlipAngleFrontLeft,
		TireSlipAngleFrontRight:              f.TireSlipAngleFrontRight,
		TireSlipAngleRearLeft:                f.TireSlipAngleRearLeft,
		TireSlipAngleRearRight:               f.TireSlipAngleRearRight,
		TireCombinedSlipFrontLeft:            f.TireCombinedSlipFrontLeft,
		TireCombinedSlipFrontRight:           f.TireCombinedSlipFrontRight,
		TireCombinedSlipRearLeft:             f.TireCombinedSlipRearLeft,
		TireCombinedSlipRearRight:            f.TireCombinedSlipRearRight,
		SuspensionTravelMetersFrontLeft:      f.SuspensionTravelMetersFrontLeft,
		SuspensionTravelMetersFrontRight:     f.SuspensionTravelMetersFrontRight,
		SuspensionTravelMetersRearLeft:       f.SuspensionTravelMetersRearLeft,
		SuspensionTravelMetersRearRight:      f.SuspensionTravelMetersRearRight,
		CarOrdinal:                           f.CarOrdinal,
		CarClass:                             f.CarClass,
		CarPerformanceIndex:                  f.CarPerformanceIndex,
		DrivetrainType:                       f.DrivetrainType,
		NumCylinders:                         f.NumCylinders,
		CarCategory:                          f.CarCategory,
		HorizonPlaceholder1:                  f.HorizonPlaceholder1,
		HorizonPlaceholder2:                  f.HorizonPlaceholder2,
		PositionX:                            f.PositionX,
		PositionY:                            f.PositionY,
		PositionZ:                            f.PositionZ,
		Speed:                                f.Speed * mps_to_kmh,
		Power:                                f.Power * watts_to_bhp,
		Torque:                               f.Torque,
		TireTempFrontLeft:                    toCelsius(f.TireTempFrontLeft),
		TireTempFrontRight:                   toCelsius(f.TireTempFrontRight),
		TireTempRearLeft:                     toCelsius(f.TireTempRearLeft),
		TireTempRearRight:                    toCelsius(f.TireTempRearRight),
		Boost:                                f.Boost,
		Fuel:                                 f.Fuel,
		DistanceTraveled:                     f.DistanceTraveled,
		BestLap:                              f.BestLap,
		LastLap:                              f.LastLap,
		CurrentLap:                           f.CurrentLap,
		CurrentRaceTime:                      f.CurrentRaceTime,
		LapNumber:                            f.LapNumber,
		RacePosition:                         f.RacePosition,
		Throttle:                             float32(f.Throttle) / 255,
		Brake:                                float32(f.Brake) / 255,
		Clutch:                               float32(f.Clutch) / 255,
		HandBrake:                            float32(f.HandBrake) / 255,
		Gear:                                 f.Gear,
		Steer:                                float32(f.Steer) / 255,
		NormalizedDrivingLine:                f.NormalizedDrivingLine,
		NormalizedAIBrakeDifference:          f.NormalizedAIBrakeDifference,
		CarAttitude:                          carAttitude,
		IsTractionLost:                       isTractionLost,
	}
	return fC
}
