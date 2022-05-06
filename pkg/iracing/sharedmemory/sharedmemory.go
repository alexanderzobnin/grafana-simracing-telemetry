//go:build windows

package sharedmemory

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"time"

	"github.com/alexeymaximov/go-bio/mmap"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"golang.org/x/sys/windows"
)

var lastTickCount int32 = IntMax

var varHeadersMap map[string]IRSDKVarHeaderDTO = map[string]IRSDKVarHeaderDTO{}

func RunSharedMemoryClient(ch chan IRacingTelemetryMap, ctrl chan string, interval time.Duration) {
	mapping, err := openMapping()
	if err != nil {
		log.DefaultLogger.Error("Error opening file mapping", "error", err)
	}
	defer closeMapping(mapping)

	hDataEvent, err := openEvent()
	if err != nil {
		log.DefaultLogger.Error("Error opening data event", "error", err)
	}

	tick := time.Tick(interval)
	for {
		select {
		case <-tick:
			header, err := readHeader(mapping)
			if err != nil {
				log.DefaultLogger.Warn("Error reading file mapping", "error", err)
				continue
			}
			data, err := readData(mapping, header, hDataEvent)
			if err != nil {
				log.DefaultLogger.Warn("Error reading file mapping", "error", err)
				continue
			}

			if data != nil {
				telemetry := convertToTelemetryMap(data)
				ch <- telemetry
			}

		case ctrlMessage := <-ctrl:
			if ctrlMessage == "stop" {
				log.DefaultLogger.Info("Stopping shared memory client")
				return
			}
		}
	}
}

func openMapping() (*mmap.Mapping, error) {
	fileSize := IRacingMemMapFileSize
	mappingNamePtr, err := windows.UTF16FromString(IRacingMemMapFileName)

	mapping, err := mmap.OpenFileMapping(INVALID_HANDLE_VALUE, 0, uintptr(fileSize), mmap.ModeReadOnly, 0, &mappingNamePtr[0])
	if err != nil {
		return nil, err
	}

	return mapping, nil
}

func openEvent() (syscall.Handle, error) {
	eventNamePtr, _ := windows.UTF16FromString(IRSDK_DATAVALIDEVENTNAME)
	hDataValidEvent, err := windows.OpenEvent(windows.SYNCHRONIZE, false, &eventNamePtr[0])
	if err != nil {
		return 0, err
	}

	return syscall.Handle(hDataValidEvent), nil
}

func waitForValidData(hDataValidEvent syscall.Handle) bool {
	event, err := syscall.WaitForSingleObject(syscall.Handle(hDataValidEvent), 100)
	if err != nil {
		return false
	}

	return event == windows.WAIT_OBJECT_0
}

func readHeader(mapping *mmap.Mapping) (*IRSDKHeader, error) {
	header := IRSDKHeader{}
	headerSize := binary.Size(header)

	b := make([]byte, headerSize)
	_, err := mapping.ReadAt(b, 0)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	b = make([]byte, header.SessionInfoLen)
	_, err = mapping.ReadAt(b, int64(header.SessionInfoOffset))
	if err != nil {
		return nil, err
	}

	// Read variables headers
	varHeaderSize := binary.Size(IRSDKVarHeader{})
	for i := 0; int32(i) < header.NumVars; i++ {
		varHeader := IRSDKVarHeader{}
		b := make([]byte, varHeaderSize)
		varHeaderOffset := header.VarHeaderOffset + int32(i*varHeaderSize)
		_, err := mapping.ReadAt(b, int64(varHeaderOffset))
		if err != nil {
			return nil, err
		}
		buf := bytes.NewReader(b)
		err = binary.Read(buf, binary.LittleEndian, &varHeader)
		if err != nil {
			return nil, err
		}

		varName := strings.Trim(string(varHeader.Name[:]), "\u0000")
		if _, ok := varHeadersMap[varName]; varName != "" && !ok {
			varHeaderDTO := varHeader.toIRSDKVarHeaderDTO()
			varHeadersMap[varName] = varHeaderDTO
		}
	}

	return &header, nil
}

func readData(mapping *mmap.Mapping, header *IRSDKHeader, hDataValidEvent syscall.Handle) ([]byte, error) {
	isValidData := waitForValidData(hDataValidEvent)
	if !isValidData {
		return nil, nil
	}

	if header.Status == 0 {
		lastTickCount = IntMax
		return nil, errors.New("client disconnected")
	}

	latest := 0
	for i := 1; int32(i) < header.NumBuf; i++ {
		if header.VarBuf[i].TickCount > header.VarBuf[latest].TickCount {
			latest = i
		}
	}

	// if newer than last recieved, than report new data
	if lastTickCount < header.VarBuf[latest].TickCount {
		// try twice to get the data out
		for i := 0; i < 2; i++ {
			curTickCount := header.VarBuf[latest].TickCount
			dataLen := header.BufLen
			dataOffset := header.VarBuf[latest].BufOffset
			b := make([]byte, dataLen)

			_, err := mapping.ReadAt(b, int64(dataOffset))
			if err != nil {
				return nil, err
			}
			if curTickCount == header.VarBuf[latest].TickCount {
				lastTickCount = curTickCount
			}
			return b, nil
		}
	} else if lastTickCount > header.VarBuf[latest].TickCount {
		lastTickCount = header.VarBuf[latest].TickCount
	} else {
		log.DefaultLogger.Debug("No data to read", "lastTickCount", lastTickCount, "TickCount", header.VarBuf[latest].TickCount)
	}

	return nil, nil
}

func closeMapping(mapping *mmap.Mapping) {
	err := mapping.Close()
	if err != nil {
		log.DefaultLogger.Warn("Error closing file mapping", "error", err)
	}
}

func convertToTelemetryMap(data []byte) IRacingTelemetryMap {
	telemetryMap := make(IRacingTelemetryMap)
	for name, varHeader := range varHeadersMap {
		switch varHeader.Type {
		case 0:
			// char
			// TODO: implement this
			break
		case 1:
			// bool
			var value bool
			valueByte := data[varHeader.Offset : varHeader.Offset+varHeader.Length]
			buf := bytes.NewReader(valueByte)
			err := binary.Read(buf, binary.LittleEndian, &value)
			if err != nil {
				log.DefaultLogger.Warn("Error converting bool value", "error", err)
				break
			}
			telemetryMap[name] = IRacingTelemetryValue{Name: name, Type: "bool", Value: value}
		case 2:
			// int
			var value int32
			valueByte := data[varHeader.Offset : varHeader.Offset+varHeader.Length]
			buf := bytes.NewReader(valueByte)
			err := binary.Read(buf, binary.LittleEndian, &value)
			if err != nil {
				log.DefaultLogger.Warn("Error converting int value", "error", err)
				break
			}
			telemetryMap[name] = IRacingTelemetryValue{Name: name, Type: "int32", Value: value}
		case 3:
			// bit field
			var value [4]byte
			valueByte := data[varHeader.Offset : varHeader.Offset+varHeader.Length]
			buf := bytes.NewReader(valueByte)
			err := binary.Read(buf, binary.LittleEndian, &value)
			if err != nil {
				log.DefaultLogger.Warn("Error converting byte value", "error", err)
				break
			}
			flagValue := binary.LittleEndian.Uint32(value[:])
			telemetryMap[name] = IRacingTelemetryValue{Name: name, Type: "uint32", Value: flagValue}
			if name == "EngineWarnings" {
				readEngineWarnings(telemetryMap, flagValue)
			}
		case 4:
			// float32
			var value float32
			valueByte := data[varHeader.Offset : varHeader.Offset+varHeader.Length]
			buf := bytes.NewReader(valueByte)
			err := binary.Read(buf, binary.LittleEndian, &value)
			if err != nil {
				log.DefaultLogger.Warn("Error converting float value", "error", err)
				break
			}

			if name == "Speed" {
				speedKmph := value * 3.6
				telemetryMap["SpeedKmh"] = IRacingTelemetryValue{Name: "SpeedKmh", Type: "float32", Value: speedKmph}
			}
			if name == "Clutch" {
				value = 1 - value
			}
			telemetryMap[name] = IRacingTelemetryValue{Name: name, Type: "float32", Value: value}
		case 5:
			// float64 (double)
			var value float64
			valueByte := data[varHeader.Offset : varHeader.Offset+varHeader.Length]
			buf := bytes.NewReader(valueByte)
			err := binary.Read(buf, binary.LittleEndian, &value)
			if err != nil {
				log.DefaultLogger.Warn("Error converting double value", "error", err)
				break
			}
			telemetryMap[name] = IRacingTelemetryValue{Name: name, Type: "float64", Value: value}
		}
	}

	return telemetryMap
}

func readEngineWarnings(telementryMap IRacingTelemetryMap, value uint32) {
	telementryMap["WaterTempWarning"] = IRacingTelemetryValue{Name: "WaterTempWarning", Type: "int", Value: 0}
	telementryMap["FuelPressureWarning"] = IRacingTelemetryValue{Name: "FuelPressureWarning", Type: "int", Value: 0}
	telementryMap["OilPressureWarning"] = IRacingTelemetryValue{Name: "OilPressureWarning", Type: "int", Value: 0}
	telementryMap["EngineStalled"] = IRacingTelemetryValue{Name: "EngineStalled", Type: "int", Value: 0}
	telementryMap["PitSpeedLimiter"] = IRacingTelemetryValue{Name: "PitSpeedLimiter", Type: "int", Value: 0}
	telementryMap["RevLimiterActive"] = IRacingTelemetryValue{Name: "RevLimiterActive", Type: "int", Value: 0}
	telementryMap["OilTempWarning"] = IRacingTelemetryValue{Name: "OilTempWarning", Type: "int", Value: 0}
	if value&1 != 0 {
		telementryMap["WaterTempWarning"] = IRacingTelemetryValue{Name: "WaterTempWarning", Type: "int", Value: 1}
	}
	if value&2 != 0 {
		telementryMap["FuelPressureWarning"] = IRacingTelemetryValue{Name: "FuelPressureWarning", Type: "int", Value: 1}
	}
	if value&4 != 0 {
		telementryMap["OilPressureWarning"] = IRacingTelemetryValue{Name: "OilPressureWarning", Type: "int", Value: 1}
	}
	if value&8 != 0 {
		telementryMap["EngineStalled"] = IRacingTelemetryValue{Name: "EngineStalled", Type: "int", Value: 1}
	}
	if value&16 != 0 {
		telementryMap["PitSpeedLimiter"] = IRacingTelemetryValue{Name: "PitSpeedLimiter", Type: "int", Value: 1}
	}
	if value&32 != 0 {
		telementryMap["RevLimiterActive"] = IRacingTelemetryValue{Name: "RevLimiterActive", Type: "int", Value: 1}
	}
	if value&64 != 0 {
		telementryMap["OilTempWarning"] = IRacingTelemetryValue{Name: "OilTempWarning", Type: "int", Value: 1}
	}
}

func TelemetryToDataFrame(tm IRacingTelemetryMap) (*data.Frame, error) {
	frame := data.NewFrame("response")
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{time.Now()}),
	)

	for _, value := range tm {
		field, err := value.TelemetryValueToField()
		if err != nil {
			log.DefaultLogger.Debug("Error converting value to field", "error", err)
		} else {
			frame.Fields = append(frame.Fields, field)
		}
	}

	return frame, nil
}

func (telemetryValue IRacingTelemetryValue) TelemetryValueToField() (*data.Field, error) {
	switch telemetryValue.Type {
	case "float32":
		value, ok := telemetryValue.Value.(float32)
		if !ok {
			return nil, errors.New(fmt.Sprintf("error converting value to %s", telemetryValue.Type))
		}
		return data.NewField(telemetryValue.Name, nil, []float32{value}), nil
	case "float64":
		value, ok := telemetryValue.Value.(float64)
		if !ok {
			return nil, errors.New(fmt.Sprintf("error converting value to %s", telemetryValue.Type))
		}
		return data.NewField(telemetryValue.Name, nil, []float64{value}), nil
	case "int32":
		value, ok := telemetryValue.Value.(int32)
		if !ok {
			return nil, errors.New(fmt.Sprintf("error converting value %v to %s", telemetryValue.Value, telemetryValue.Type))
		}
		return data.NewField(telemetryValue.Name, nil, []int32{value}), nil
	case "uint32":
		value, ok := telemetryValue.Value.(uint32)
		if !ok {
			return nil, errors.New(fmt.Sprintf("error converting value %v to %s", telemetryValue.Value, telemetryValue.Type))
		}
		return data.NewField(telemetryValue.Name, nil, []uint32{value}), nil
	case "int":
		value, ok := telemetryValue.Value.(int)
		if !ok {
			return nil, errors.New(fmt.Sprintf("error converting value %v to %s", telemetryValue.Value, telemetryValue.Type))
		}
		return data.NewField(telemetryValue.Name, nil, []int32{int32(value)}), nil
	case "bool":
		value, ok := telemetryValue.Value.(bool)
		if !ok {
			return nil, errors.New(fmt.Sprintf("error converting value to %s", telemetryValue.Type))
		}
		valueInt := 0
		if value {
			valueInt = 1
		}
		return data.NewField(telemetryValue.Name, nil, []int32{int32(valueInt)}), nil
	default:
		return nil, errors.New(fmt.Sprintf("not implemented type %s", telemetryValue.Type))
	}
}

func readValueAsDataField(dataByte []byte, name string) (*data.Field, error) {
	varHeader, ok := varHeadersMap[name]
	if !ok {
		return nil, errors.New("field not found")
	}

	switch varHeader.Type {
	case 0:
		// char
		return nil, errors.New("not implemented")
	case 1:
		// bool
		var value bool
		valueByte := dataByte[varHeader.Offset : varHeader.Offset+varHeader.Length]
		buf := bytes.NewReader(valueByte)
		err := binary.Read(buf, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		return data.NewField(name, nil, []bool{value}), nil
	case 2:
		// int
		var value int32
		valueByte := dataByte[varHeader.Offset : varHeader.Offset+varHeader.Length]
		buf := bytes.NewReader(valueByte)
		err := binary.Read(buf, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		return data.NewField(name, nil, []int32{value}), nil
	case 3:
		// bit field
		return nil, errors.New("not implemented")
	case 4:
		// float32
		var value float32
		valueByte := dataByte[varHeader.Offset : varHeader.Offset+varHeader.Length]
		buf := bytes.NewReader(valueByte)
		err := binary.Read(buf, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		return data.NewField(name, nil, []float32{value}), nil
	case 5:
		// float64 (double)
		var value float64
		valueByte := dataByte[varHeader.Offset : varHeader.Offset+varHeader.Length]
		buf := bytes.NewReader(valueByte)
		err := binary.Read(buf, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		return data.NewField(name, nil, []float64{value}), nil
	}

	return nil, errors.New("field has unknown type")
}
