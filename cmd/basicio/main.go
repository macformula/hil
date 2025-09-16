package main

import (
	"fmt"
	"io"
	"net"
	"time"

	signals "github.com/macformula/hil/cmd/basicio/signals"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	_port = 8080
)

type BasicIo struct {
	l *zap.Logger
}

const ecu_name string = "DemoProject"
const button_name string = "IndicatorButton"
const led_name string = "IndicatorLed"

func addr(s string) *string {
	return &s
}

func button() *signals.Signal {
	return &signals.Signal{
		EcuName:    addr(ecu_name),
		SignalName: addr(button_name),
		Type:       signals.SignalType_SIGNAL_TYPE_DIGITAL.Enum(),
		Direction:  signals.SignalDirection_SIGNAL_DIRECTION_OUTPUT.Enum(),
	}
}

func led() *signals.Signal {
	return &signals.Signal{
		EcuName:    addr(ecu_name),
		SignalName: addr(led_name),
		Type:       signals.SignalType_SIGNAL_TYPE_DIGITAL.Enum(),
		Direction:  signals.SignalDirection_SIGNAL_DIRECTION_INPUT.Enum(),
	}
}

func (c *BasicIo) ReadButtonValue(conn net.Conn) (bool, error) {
	readButtonRequest := signals.Request{
		Signal:  button(),
		Request: &signals.Request_Read{},
	}

	bytes, err := proto.Marshal(&readButtonRequest)
	if err != nil {
		c.l.Error(fmt.Sprintf("Error serializing request: %v", err))
	}

	// Send data to the server
	_, err = conn.Write(bytes)
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
			return false, fmt.Errorf("read error: %s", err)
		}
	}

	var response signals.Response
	err = proto.Unmarshal(buffer, &response)
	if err != nil {
		return false, fmt.Errorf("Failed to deserialize response")
	}

	if !*response.Ok {
		return false, fmt.Errorf("Request was not ok!")
	}

	read := response.GetRead()
	if read == nil {
		return false, fmt.Errorf("Invalid response type")
	}

	if v, ok := read.Value.GetValue().(*signals.SignalValue_Digital); ok {
		return v.Digital, nil
	} else {
		return false, fmt.Errorf("Missing digital value")
	}
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

		var response signals.Response
		err = proto.Unmarshal(buffer, &response)

		if err != nil {
			c.l.Error("Failed to deserialize Response")
			continue
		}

		c.l.Info(fmt.Sprintf("received response ok (%t) errorString (%s) level (%t) voltage (%f)", response.Ok, *response.Error, response.GetRead().GetValue().GetDigital(), response.GetRead().GetValue().GetAnalogVoltage()))
	}
}

func main() {
	fmt.Printf("starting")
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

	for {
		level, err := firmware.ReadButtonValue(conn)
		if err != nil {
			fmt.Printf("Read button value (likely timeout error)")
			return
		}
		fmt.Printf("Read indicator button is %t", level)

		setLedRequest := signals.Request{
			Signal: led(),
			Request: &signals.Request_Set{
				Set: &signals.SetRequest{
					Value: &signals.SignalValue{
						Value: &signals.SignalValue_Digital{
							Digital: level,
						},
					}},
			},
		}

		bytes, err := proto.Marshal(&setLedRequest)
		if err != nil {
			fmt.Println("Error sending data:", err)
		}

		// Send data to the server
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("Error sending data:", err)
		}
		fmt.Printf("Set indicator led to %t", level)
		time.Sleep(50 * time.Millisecond)
	}
}
