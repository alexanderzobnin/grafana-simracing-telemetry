package acc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

var server = "192.168.88.209"
var port = "9000"
var connectionPassword = "acc"
var commandPassword = ""
var displayName = "Grafana"
var msUpdateInterval int32 = 250

func RunClient(errCh chan error) {
	serverAddr := fmt.Sprintf("%s:%s", server, port)
	s, err := net.ResolveUDPAddr("udp4", serverAddr)
	if err != nil {
		errCh <- err
		return
	}

	conn, err := net.DialUDP("udp", nil, s)
	if err != nil {
		errCh <- err
		return
	}
	log.DefaultLogger.Info("Connecting UDP telemetry server")
	defer conn.Close()

	err = RequestConnection(conn, displayName, connectionPassword, commandPassword, msUpdateInterval)
	if err != nil {
		log.DefaultLogger.Error("Error connecting UDP server", "error", err)
	}

	buffer := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			errCh <- err
			return
		}
		log.DefaultLogger.Debug("Received message", "n", n)

		messageType := uint8(buffer[0])
		log.DefaultLogger.Debug("Message type", "type", messageType)

		if messageType == REGISTRATION_RESULT {
			rr := HandleRegistrationResult(buffer)
			log.DefaultLogger.Debug("Registration", "result", rr)
		} else if messageType == REALTIME_CAR_UPDATE {
			cu := HandleCarUpdate(buffer)
			log.DefaultLogger.Debug("Car update", "result", cu)
		}
	}
}

func RequestConnection(conn *net.UDPConn, displayName, connectionPassword, commandPassword string, msUpdateInterval int32) error {
	bs := make([]byte, 0)
	buf := bytes.NewBuffer(bs)

	binary.Write(buf, binary.LittleEndian, int8(1))
	binary.Write(buf, binary.LittleEndian, int8(4))

	strLen := uint16(len(displayName))
	binary.Write(buf, binary.LittleEndian, strLen)
	fmt.Fprint(buf, displayName)

	strLen = uint16(len(connectionPassword))
	binary.Write(buf, binary.LittleEndian, strLen)
	fmt.Fprint(buf, connectionPassword)

	binary.Write(buf, binary.LittleEndian, msUpdateInterval)

	strLen = uint16(len(commandPassword))
	binary.Write(buf, binary.LittleEndian, strLen)
	fmt.Fprint(buf, commandPassword)

	n, err := conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	log.DefaultLogger.Debug("Send connection request to server", "n", n)

	return nil
}

func Disconnect(conn *net.UDPConn) error {
	bs := make([]byte, 0)
	buf := bytes.NewBuffer(bs)

	binary.Write(buf, binary.LittleEndian, int8(9))

	n, err := conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	log.DefaultLogger.Debug("Send disconnect request to server", "n", n)

	return nil
}
