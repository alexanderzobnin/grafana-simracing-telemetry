package sharedmemory

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

// IRacingTelemetry represents
type IRacingTelemetry struct {
	PacketId int32
	Gas      float32
}

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
	BufOffset int32    // offset from header
	TickCount int32    // used to detect changes in data
	Pad       [2]int32 // (16 byte align)
}
