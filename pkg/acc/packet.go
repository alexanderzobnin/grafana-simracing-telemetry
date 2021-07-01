package acc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	REGISTRATION_RESULT = 1
	REALTIME_UPDATE     = 2
	REALTIME_CAR_UPDATE = 3
)

type RegistrationResult struct {
	ConnectionID int32
	Success      bool
	IsReadOnly   bool
	ErrMessage   string
}

type CarUpdate struct {
	CarIndex    uint16
	DriverIndex uint16
	DriverCount uint8
	Gear        int8
	WorldPosX   float32
	WorldPosY   float32
	Yaw         float32
	CarLocation uint8
	Kmh         uint16
}

func HandleRegistrationResult(b []byte) RegistrationResult {
	rr := RegistrationResult{}
	buf := bytes.NewReader(b)
	var messageType uint8 = 0
	binary.Read(buf, binary.LittleEndian, &messageType)

	binary.Read(buf, binary.LittleEndian, &rr.ConnectionID)
	binary.Read(buf, binary.LittleEndian, &rr.Success)
	binary.Read(buf, binary.LittleEndian, &rr.IsReadOnly)
	rr.ErrMessage = readString(buf)

	return rr
}

func HandleCarUpdate(b []byte) CarUpdate {
	cu := CarUpdate{}
	buf := bytes.NewReader(b)
	var messageType uint8 = 0
	binary.Read(buf, binary.LittleEndian, &messageType)

	binary.Read(buf, binary.LittleEndian, &cu)
	cu.Gear = cu.Gear - 1

	return cu
}

func writeString(buf io.Writer, str string) {
	strLen := uint16(len(str))
	binary.Write(buf, binary.LittleEndian, strLen)
	fmt.Fprint(buf, str)
}

func readString(buf io.Reader) string {
	var strLen uint16
	binary.Read(buf, binary.LittleEndian, strLen)
	str := make([]byte, strLen)
	buf.Read(str)
	return string(str)
}
