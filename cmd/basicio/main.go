package main

import (
	"fmt"
	"io"
	"net"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	signals "github.com/macformula/hil/cmd/basicio/signals"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_port = 8080
)

type BasicIo struct {
	l *zap.Logger
}

func (c *BasicIo) ReadButtonValue(conn net.Conn) (bool, error) {
	readButtonRequest := serializeReadRequest("DemoProject", "IndicatorButton", signals.SIGNAL_TYPEDIGITAL, signals.SIGNAL_DIRECTIONINPUT)

	// Send data to the server
	_, err := conn.Write(readButtonRequest)
	if err != nil {
		c.l.Error(fmt.Sprintf("Error sending data: %s", err))
	}
	c.l.Info("Send read button request")
	level, err := c.WaitForResponse(conn)
	return level, err
}

func (c *BasicIo) WaitForResponse(conn net.Conn) (bool, error) {
	// Load server's response into buffer
	buffer := make([]byte, 2024)
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			c.l.Error(fmt.Sprintf("read error: %s", err))
			return false, errors.Wrap(err, "")
		}
	}

	// Deserilize server's response
	response := signals.GetRootAsResponse(buffer, 0)
	unionTable := new(flatbuffers.Table)
	if response.Response(unionTable) {
		ok, errorString, level, voltage := deserializeReadResponse(unionTable)
		c.l.Info(fmt.Sprintf("recieved response ok (%t) errorString (%s) level (%t) voltage (%f)", ok, errorString, level, voltage))
		return level, nil
	}
	return false, errors.Errorf("Read reponse could not fit in table.")
}

func (c *BasicIo) Listen(conn net.Conn) {
	for {
		buffer := make([]byte, 2024)
		_, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				c.l.Error(fmt.Sprintf("read error: %s", err))
			}
		}
		response := signals.GetRootAsResponse(buffer, 0)
		unionTable := new(flatbuffers.Table)
		if response.Response(unionTable) {
			ok, errorString, level, voltage := deserializeReadResponse(unionTable)
			c.l.Info(fmt.Sprintf("recieved response ok (%t) errorString (%s) level (%t) voltage (%f)", ok, errorString, level, voltage))
		}
	}
}

func main() {
	loggerConfig := zap.NewDevelopmentConfig()
	logger, _ := loggerConfig.Build()

	time.Sleep(2 * time.Second)

	// Connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", _port))
	if err != nil {
		logger.Info(fmt.Sprintf("Error connecting to server: %e", err))
	}
	defer conn.Close()

	fmt.Println("Connected to sil server")

	firmware := BasicIo{
		l: logger,
	}

	registerIndicatorLed := serializeRegisterRequest("DemoProject", "IndicatorLed", signals.SIGNAL_TYPEDIGITAL, signals.SIGNAL_DIRECTIONOUTPUT)
	registerIndicatorButton := serializeRegisterRequest("DemoProject", "IndicatorButton", signals.SIGNAL_TYPEDIGITAL, signals.SIGNAL_DIRECTIONINPUT)

	_, err = conn.Write(registerIndicatorLed)
	if err != nil {
		firmware.l.Error(fmt.Sprintf("Error sending data: %s", err))
	}
	firmware.l.Info("Registered indicator led")

	_, err = conn.Write(registerIndicatorButton)
	if err != nil {
		firmware.l.Error(fmt.Sprintf("Error sending data: %s", err))
	}
	firmware.l.Info("Registered indicator button")

	for {
		level, err := firmware.ReadButtonValue(conn)
		if err != nil {
			firmware.l.Error("Read button value (likely timeout error)")
		}
		firmware.l.Info(fmt.Sprintf("Read indicator button is %t", level))

		setLedRequest := serializeSetRequest("DemoProject", "IndicatorLed", signals.SIGNAL_TYPEDIGITAL, 0.0, level)
		// Send data to the server
		_, err = conn.Write(setLedRequest)
		if err != nil {
			fmt.Println("Error sending data:", err)
		}
		firmware.l.Info(fmt.Sprintf("Set indicator led to %t", level))

		time.Sleep(100 * time.Millisecond)
	}
}
