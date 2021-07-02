package sharedmemory

import (
	"bytes"
	"encoding/binary"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"golang.org/x/sys/windows"
	"time"

	"github.com/alexeymaximov/go-bio/mmap"
)

const (
	INVALID_HANDLE_VALUE uintptr = 0x80000000
	ACCPageFilePhysics           = "Local\\acpmf_physics"
	ACCPageFileGraphic           = "Local\\acpmf_graphics"
	ACCPageFileStatic            = "Local\\acpmf_static"
)

func RunSharedMemoryClient(ch chan SPageFilePhysics, ctrl chan string, interval time.Duration) {
	mapping, err := openPhysicsMapping()
	if err != nil {
		log.DefaultLogger.Error("Error opening mapping", "error", err)
	}
	defer closeMapping(mapping)

	tick := time.Tick(interval)
	for {
		select {
		case <-tick:
			physics, err := readPhysics(mapping)
			if err != nil {
				log.DefaultLogger.Warn("Error reading file mapping", "error", err)
			}
			ch <- *physics
			break
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

func openPhysicsMapping() (*mmap.Mapping, error) {
	fileSize := binary.Size(SPageFilePhysics{})
	mappingNamePtr, err := windows.UTF16FromString(ACCPageFilePhysics)

	mapping, err := mmap.OpenFileMapping(INVALID_HANDLE_VALUE, 0, uintptr(fileSize), mmap.ModeReadWrite, 0, &mappingNamePtr[0])
	if err != nil {
		return nil, err
	}

	return mapping, nil
}

func closeMapping(mapping *mmap.Mapping) {
	err := mapping.Close()
	if err != nil {
		log.DefaultLogger.Error("Error closing file mapping", "error", err)
	}
}
