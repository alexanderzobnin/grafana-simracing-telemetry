package sharedmemory

import (
	"bytes"
	"encoding/binary"
	"github.com/alexeymaximov/go-bio/mmap"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"golang.org/x/sys/windows"
	"time"
)

func RunSharedMemoryClient(ch chan IRacingTelemetry, ctrl chan string, interval time.Duration) {
	mapping, err := openMapping()
	if err != nil {
		log.DefaultLogger.Error("Error opening file mapping", "error", err)
	}
	defer closeMapping(mapping)

	tick := time.Tick(interval)
	for {
		select {
		case <-tick:
			//data, err := readTelemetry(mapping)
			header, err := readHeader(mapping)
			log.DefaultLogger.Debug("Header", "header", header)
			if err != nil {
				log.DefaultLogger.Warn("Error reading file mapping", "error", err)
			}

			//ch <- *data

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

	mapping, err := mmap.OpenFileMapping(INVALID_HANDLE_VALUE, 0, uintptr(fileSize), mmap.ModeReadWrite, 0, &mappingNamePtr[0])
	if err != nil {
		return nil, err
	}

	return mapping, nil
}

func readHeader(mapping *mmap.Mapping) (*IRSDKHeader, error) {
	data := IRSDKHeader{}
	fileSize := binary.Size(data)

	b := make([]byte, fileSize)
	_, err := mapping.ReadAt(b, 0)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func readTelemetry(mapping *mmap.Mapping) (*IRacingTelemetry, error) {
	data := IRacingTelemetry{}
	//fileSize := binary.Size(data)
	fileSize := IRacingMemMapFileSize

	b := make([]byte, fileSize)
	_, err := mapping.ReadAt(b, 0)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func closeMapping(mapping *mmap.Mapping) {
	err := mapping.Close()
	if err != nil {
		log.DefaultLogger.Warn("Error closing file mapping", "error", err)
	}
}
