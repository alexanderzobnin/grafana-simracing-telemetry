//go:build windows

package sharedmemory

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"golang.org/x/sys/windows"

	"github.com/alexeymaximov/go-bio/mmap"
)

const (
	INVALID_HANDLE_VALUE uintptr = 0x80000000
	ACCPageFilePhysics           = "Local\\acpmf_physics"
	ACCPageFileGraphic           = "Local\\acpmf_graphics"
	ACCPageFileStatic            = "Local\\acpmf_static"
)

func RunSharedMemoryClient(ch chan ACCTelemetry, ctrl chan string, interval time.Duration) {
	mPhysics, err := openPhysicsMapping()
	if err != nil {
		log.DefaultLogger.Error("Error opening physics file mapping", "error", err)
	}
	defer closeMapping(mPhysics)

	mGraphic, err := openGraphicMapping()
	if err != nil {
		log.DefaultLogger.Error("Error opening graphic file mapping", "error", err)
	}
	defer closeMapping(mGraphic)

	tick := time.Tick(interval)
	for {
		select {
		case <-tick:
			physics, err := readPhysics(mPhysics)
			if err != nil {
				log.DefaultLogger.Warn("Error reading physics file mapping", "error", err)
			}

			graphic, err := readGraphic(mGraphic)
			if err != nil {
				log.DefaultLogger.Warn("Error reading graphic file mapping", "error", err)
			}

			accTelemetry := ACCTelemetry{
				SPageFileGraphic: *graphic,
				SPageFilePhysics: *physics,
			}
			ch <- accTelemetry

		case ctrlMessage := <-ctrl:
			if ctrlMessage == "stop" {
				log.DefaultLogger.Info("Stopping shared memory client")
				return
			}
		}
	}
}

func readPhysics(mapping *mmap.Mapping) (*SPageFilePhysics, error) {
	physics := SPageFilePhysics{}
	fileSize := binary.Size(physics)

	b := make([]byte, fileSize)
	_, err := mapping.ReadAt(b, 0)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &physics)
	if err != nil {
		return nil, err
	}

	return &physics, nil
}

func readGraphic(mapping *mmap.Mapping) (*SPageFileGraphic, error) {
	graphic := SPageFileGraphic{}
	fileSize := binary.Size(graphic)

	b := make([]byte, fileSize)
	_, err := mapping.ReadAt(b, 0)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &graphic)
	if err != nil {
		return nil, err
	}

	return &graphic, nil
}

func openPhysicsMapping() (*mmap.Mapping, error) {
	fileSize := binary.Size(SPageFilePhysics{})
	mappingNamePtr, err := windows.UTF16FromString(ACCPageFilePhysics)

	mapping, err := mmap.OpenFileMapping(INVALID_HANDLE_VALUE, 0, uintptr(fileSize), mmap.ModeReadWrite, 0, &mappingNamePtr[0])
	if err != nil {
		return nil, err
	}

	return mapping, nil
}

func openGraphicMapping() (*mmap.Mapping, error) {
	fileSize := binary.Size(SPageFileGraphic{})
	mappingNamePtr, err := windows.UTF16FromString(ACCPageFileGraphic)

	mapping, err := mmap.OpenFileMapping(INVALID_HANDLE_VALUE, 0, uintptr(fileSize), mmap.ModeReadWrite, 0, &mappingNamePtr[0])
	if err != nil {
		return nil, err
	}

	return mapping, nil
}

func closeMapping(mapping *mmap.Mapping) {
	err := mapping.Close()
	if err != nil {
		log.DefaultLogger.Warn("Error closing file mapping", "error", err)
	}
}
