package sharedmemory

import "strings"

const (
	INVALID_HANDLE_VALUE     uintptr = 0x80000000
	IRSDK_DATAVALIDEVENTNAME         = "Local\\IRSDKDataValidEvent"

	IRacingMemMapFileName = "Local\\IRSDKMemMapFileName"
	IRacingMemMapFileSize = 1164 * 1024

	IRSDK_MAX_STRING = 32
	IRSDK_MAX_DESC   = 64
	IRSDK_MAX_BUFS   = 4

	IntMax = 2147483647
)

type IRSDKVarHeader struct {
	Type   int32 // irsdk_VarType
	Offset int32 // offset fron start of buffer row
	Count  int32 // number of entrys (array) so length in bytes would be irsdk_VarTypeBytes[type] * count

	CountAsTime bool
	Pad         [3]byte // (16 byte align)

	Name [IRSDK_MAX_STRING]byte
	Desc [IRSDK_MAX_DESC]byte
	Unit [IRSDK_MAX_STRING]byte // something like "kg/m^2"
}

type IRSDKVarHeaderDTO struct {
	Type   int32 // irsdk_VarType
	Offset int32 // offset fron start of buffer row
	Count  int32 // number of entrys (array) so length in bytes would be irsdk_VarTypeBytes[type] * count

	CountAsTime bool

	Name string
	Desc string
	Unit string // something like "kg/m^2"
}

func (vh IRSDKVarHeader) toIRSDKVarHeaderDTO() IRSDKVarHeaderDTO {
	vhDTO := IRSDKVarHeaderDTO{
		Type:        vh.Type,
		Offset:      vh.Offset,
		Count:       vh.Count,
		CountAsTime: vh.CountAsTime,
		Name:        strings.Trim(string(vh.Name[:]), "\u0000"),
		Desc:        strings.Trim(string(vh.Desc[:]), "\u0000"),
		Unit:        strings.Trim(string(vh.Unit[:]), "\u0000"),
	}

	return vhDTO
}

type IRSDKHeader struct {
	Ver      int32 // this api header version, see IRSDK_VER
	Status   int32 // bitfield using irsdk_StatusField
	TickRate int32 // ticks per second (60 or 360 etc)

	// session information, updated periodicaly
	SessionInfoUpdate int32 // Incremented when session info changes
	SessionInfoLen    int32 // Length in bytes of session info string
	SessionInfoOffset int32 // Session info, encoded in YAML format

	// State data, output at tickRate

	NumVars         int32 // length of arra pointed to by varHeaderOffset
	VarHeaderOffset int32 // offset to irsdk_varHeader[numVars] array, Describes the variables received in varBuf

	NumBuf int32                       // <= IRSDK_MAX_BUFS (3 for now)
	BufLen int32                       // length in bytes for one line
	Pad1   [2]int32                    // (16 byte align)
	VarBuf [IRSDK_MAX_BUFS]IRSDKVarBuf // buffers of data being written to
}

type IRSDKVarBuf struct {
	TickCount int32    // used to detect changes in data
	BufOffset int32    // offset from header
	Pad       [2]int32 // (16 byte align)
}

type IRacingTelemetry struct {
	SessionTime                     float64
	SessionTick                     int32
	SessionNum                      int32
	SessionState                    int32
	SessionUniqueID                 int32
	SessionFlags                    [4]byte
	SessionTimeRemain               float64
	SessionLapsRemain               int32
	SessionLapsRemainEx             int32
	SessionTimeTotal                float64
	SessionLapsTotal                int32
	SessionTimeOfDay                float32
	RadioTransmitCarIdx             int32
	RadioTransmitRadioIdx           int32
	RadioTransmitFrequencyIdx       int32
	DisplayUnits                    int32
	DriverMarker                    bool
	PushToPass                      bool
	ManualBoost                     bool
	ManualNoBoost                   bool
	IsOnTrack                       bool
	IsReplayPlaying                 bool
	ReplayFrameNum                  int32
	ReplayFrameNumEnd               int32
	IsDiskLoggingEnabled            bool
	IsDiskLoggingActive             bool
	FrameRate                       float32
	CpuUsageFG                      float32
	GpuUsage                        float32
	ChanAvgLatency                  float32
	ChanLatency                     float32
	ChanQuality                     float32
	ChanPartnerQuality              float32
	CpuUsageBG                      float32
	ChanClockSkew                   float32
	MemPageFaultSec                 float32
	PlayerCarPosition               int32
	PlayerCarClassPosition          int32
	PlayerCarClass                  int32
	PlayerTrackSurface              int32
	PlayerTrackSurfaceMaterial      int32
	PlayerCarIdx                    int32
	PlayerCarTeamIncidentCount      int32
	PlayerCarMyIncidentCount        int32
	PlayerCarDriverIncidentCount    int32
	PlayerCarWeightPenalty          float32
	PlayerCarPowerAdjust            float32
	PlayerCarDryTireSetLimit        int32
	PlayerCarTowTime                float32
	PlayerCarInPitStall             bool
	PlayerCarPitSvStatus            int32
	PlayerTireCompound              int32
	PlayerFastRepairsUsed           int32
	CarIdxLap                       int32
	CarIdxLapCompleted              int32
	CarIdxLapDistPct                float32
	CarIdxTrackSurface              int32
	CarIdxTrackSurfaceMaterial      int32
	CarIdxOnPitRoad                 bool
	CarIdxPosition                  int32
	CarIdxClassPosition             int32
	CarIdxClass                     int32
	CarIdxF2Time                    float32
	CarIdxEstTime                   float32
	CarIdxLastLapTime               float32
	CarIdxBestLapTime               float32
	CarIdxBestLapNum                int32
	CarIdxTireCompound              int32
	CarIdxQualTireCompound          int32
	CarIdxQualTireCompoundLocked    bool
	CarIdxFastRepairsUsed           int32
	PaceMode                        int32
	CarIdxPaceLine                  int32
	CarIdxPaceRow                   int32
	CarIdxPaceFlags                 int32
	OnPitRoad                       bool
	CarIdxSteer                     float32
	CarIdxRPM                       float32
	CarIdxGear                      int32
	SteeringWheelAngle              float32
	Throttle                        float32
	Brake                           float32
	Clutch                          float32
	Gear                            int32
	RPM                             float32
	Lap                             int32
	LapCompleted                    int32
	LapDist                         float32
	LapDistPct                      float32
	RaceLaps                        int32
	LapBestLap                      int32
	LapBestLapTime                  float32
	LapLastLapTime                  float32
	LapCurrentLapTime               float32
	LapLasNLapSeq                   int32
	LapLastNLapTime                 float32
	LapBestNLapLap                  int32
	LapBestNLapTime                 float32
	LapDeltaToBestLap               float32
	LapDeltaToBestLap_DD            float32
	LapDeltaToBestLap_OK            bool
	LapDeltaToOptimalLap            float32
	LapDeltaToOptimalLap_DD         float32
	LapDeltaToOptimalLap_OK         bool
	LapDeltaToSessionBestLap        float32
	LapDeltaToSessionBestLap_DD     float32
	LapDeltaToSessionBestLap_OK     bool
	LapDeltaToSessionOptimalLap     float32
	LapDeltaToSessionOptimalLap_DD  float32
	LapDeltaToSessionOptimalLap_OK  bool
	LapDeltaToSessionLastlLap       float32
	LapDeltaToSessionLastlLap_DD    float32
	LapDeltaToSessionLastlLap_OK    bool
	Speed                           float32
	Yaw                             float32
	YawNorth                        float32
	Pitch                           float32
	Roll                            float32
	EnterExitReset                  int32
	TrackTemp                       float32
	TrackTempCrew                   float32
	AirTemp                         float32
	WeatherType                     int32
	Skies                           int32
	AirDensity                      float32
	AirPressure                     float32
	WindVel                         float32
	WindDir                         float32
	RelativeHumidity                float32
	FogLevel                        float32
	DCLapStatus                     int32
	DCDriversSoFar                  int32
	OkToReloadTextures              bool
	LoadNumTextures                 bool
	CarLeftRight                    [4]byte
	PitsOpen                        bool
	VidCapEnabled                   bool
	VidCapActive                    bool
	PitRepairLeft                   float32
	PitOptRepairLeft                float32
	PitstopActive                   bool
	FastRepairUsed                  int32
	FastRepairAvailable             int32
	LFTiresUsed                     int32
	RFTiresUsed                     int32
	LRTiresUsed                     int32
	RRTiresUsed                     int32
	LeftTireSetsUsed                int32
	RightTireSetsUsed               int32
	FrontTireSetsUsed               int32
	RearTireSetsUsed                int32
	TireSetsUsed                    int32
	LFTiresAvailable                int32
	RFTiresAvailable                int32
	LRTiresAvailable                int32
	RRTiresAvailable                int32
	LeftTireSetsAvailable           int32
	RightTireSetsAvailable          int32
	FrontTireSetsAvailable          int32
	RearTireSetsAvailable           int32
	TireSetsAvailable               int32
	CamCarIdx                       int32
	CamCameraNumber                 int32
	CamGroupNumber                  int32
	CamCameraState                  [4]byte
	IsOnTrackCar                    bool
	IsInGarage                      bool
	SteeringWheelPctTorque          float32
	SteeringWheelPctTorqueSign      float32
	SteeringWheelPctTorqueSignStops float32
	SteeringWheelPctDamper          float32
	SteeringWheelAngleMax           float32
	SteeringWheelLimiter            float32
	ShiftIndicatorPct               float32
	ShiftPowerPct                   float32
	ShiftGrindRPM                   float32
	ThrottleRaw                     float32
	BrakeRaw                        float32
	HandbrakeRaw                    float32
	SteeringWheelPeakForceNm        float32
	SteeringWheelMaxForceNm         float32
	SteeringWheelUseLinear          bool
	BrakeABSactive                  bool
	EngineWarnings                  [4]byte
	FuelLevel                       float32
	FuelLevelPct                    float32
	PitSvFlags                      [4]byte
	PitSvLFP                        float32
	PitSvRFP                        float32
	PitSvLRP                        float32
	PitSvRRP                        float32
	PitSvFuel                       float32
	PitSvTireCompound               int32
	CarIdxP2P_Status                bool
	CarIdxP2P_Count                 int32
	ReplayPlaySpeed                 int32
	ReplayPlaySlowMotion            bool
	ReplaySessionTime               float64
	ReplaySessionNum                int32
	TireLF_RumblePitch              float32
	TireRF_RumblePitch              float32
	TireLR_RumblePitch              float32
	TireRR_RumblePitch              float32
	SteeringWheelTorque_ST          float32
	SteeringWheelTorque             float32
	VelocityZ_ST                    float32
	VelocityY_ST                    float32
	VelocityX_ST                    float32
	VelocityZ                       float32
	VelocityY                       float32
	VelocityX                       float32
	YawRate_ST                      float32
	PitchRate_ST                    float32
	RollRate_ST                     float32
	YawRate                         float32
	PitchRate                       float32
	RollRate                        float32
	VertAccel_ST                    float32
	LatAccel_ST                     float32
	LongAccel_ST                    float32
	VertAccel                       float32
	LatAccel                        float32
	LongAccel                       float32
	dcStarter                       bool
	dcToggleWindshieldWipers        bool
	dcTriggerWindshieldWipers       bool
	dpRFTireChange                  float32
	dpLFTireChange                  float32
	dpRRTireChange                  float32
	dpLRTireChange                  float32
	dpFuelFill                      float32
	dpWindshieldTearoff             float32
	dpFuelAddKg                     float32
	dpFastRepair                    float32
	dcBrakeBias                     float32
	dcLaunchRPM                     float32
	dpLFTireColdPress               float32
	dpRFTireColdPress               float32
	dpLRTireColdPress               float32
	dpRRTireColdPress               float32
	WaterTemp                       float32
	WaterLevel                      float32
	FuelPress                       float32
	FuelUsePerHour                  float32
	OilTemp                         float32
	OilPress                        float32
	OilLevel                        float32
	Voltage                         float32
	ManifoldPress                   float32
	RFbrakeLinePress                float32
	RFcoldPressure                  float32
	RFtempCL                        float32
	RFtempCM                        float32
	RFtempCR                        float32
	RFwearL                         float32
	RFwearM                         float32
	RFwearR                         float32
	LFbrakeLinePress                float32
	LFcoldPressure                  float32
	LFtempCL                        float32
	LFtempCM                        float32
	LFtempCR                        float32
	LFwearL                         float32
	LFwearM                         float32
	LFwearR                         float32
	RRbrakeLinePress                float32
	RRcoldPressure                  float32
	RRtempCL                        float32
	RRtempCM                        float32
	RRtempCR                        float32
	RRwearL                         float32
	RRwearM                         float32
	RRwearR                         float32
	LRbrakeLinePress                float32
	LRcoldPressure                  float32
	LRtempCL                        float32
	LRtempCM                        float32
	LRtempCR                        float32
	LRwearL                         float32
	LRwearM                         float32
	LRwearR                         float32
	RRshockDefl                     float32
	RRshockDefl_ST                  float32
	RRshockVel                      float32
	RRshockVel_ST                   float32
	LRshockDefl                     float32
	LRshockDefl_ST                  float32
	LRshockVel                      float32
	LRshockVel_ST                   float32
	RFshockDefl                     float32
	RFshockDefl_ST                  float32
	RFshockVel                      float32
	RFshockVel_ST                   float32
	LFshockDefl                     float32
	LFshockDefl_ST                  float32
	LFshockVel                      float32
	LFshockVel_ST                   float32
}
