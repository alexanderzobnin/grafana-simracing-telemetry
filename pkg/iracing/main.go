//go:build windows

package main

import (
	"fmt"
	"os"
	"time"

	iracing "github.com/grafana/simracing-telemetry-datasource/pkg/iracing/sharedmemory"
)

func main() {
	fmt.Println("iRacing telemetry reader")

	iracingTelemetryChan := make(chan iracing.IRacingTelemetry)
	iracingCtrlChan := make(chan string)

	updateInterval := time.Millisecond * 16
	go iracing.RunSharedMemoryClient(iracingTelemetryChan, iracingCtrlChan, updateInterval)

	for {
		select {
		case <-iracingTelemetryChan:
			fmt.Println("reading data")
		}
	}

	os.Exit(1)
}
