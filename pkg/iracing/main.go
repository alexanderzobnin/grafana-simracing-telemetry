package main

import (
	"fmt"
	iracing "github.com/grafana/grafana-starter-datasource-backend/pkg/iracing/sharedmemory"
	"os"
	"time"
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
