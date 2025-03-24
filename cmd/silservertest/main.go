package main

import (
	"context"
	"fmt"
	"net"
	"time"

	signals "github.com/macformula/hil/cmd/silservertest/signals"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/macformula/hil/iocontrol/sil"
	"go.uber.org/zap"
)

type Client struct {
	addr string
	l    *zap.Logger
}

func (c *Client) Write() {
	// Connect to the server
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		c.l.Info(fmt.Sprintf("Error connecting to server: %e", err))
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	reg_req := serializeRegisterRequest("lv", "raspi_en", signals.SIGNAL_TYPEDIGITAL, signals.SIGNAL_DIRECTIONINPUT)
	read_req1 := serializeReadRequest("lv", "raspi_en", signals.SIGNAL_TYPEDIGITAL)
	set_req := serializeSetRequest("lv", "raspi_en", signals.SIGNAL_TYPEDIGITAL, 0.0, true)
	read_req2 := serializeReadRequest("lv", "raspi_en", signals.SIGNAL_TYPEDIGITAL)
	// Send data to the server
	_, err = conn.Write(reg_req)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	time.Sleep(2 * time.Second)

	_, err = conn.Write(read_req1)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	time.Sleep(2 * time.Second)

	_, err = conn.Write(set_req)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	time.Sleep(2 * time.Second)

	_, err = conn.Write(read_req2)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
}

func main() {
	ctx := context.Background()

	loggerConfig := zap.NewDevelopmentConfig()
	logger, _ := loggerConfig.Build()

	controller := sil.NewFbController(12345, logger)

	go controller.Open(ctx)

	client := Client{
		addr: "localhost:12345",
		l:    logger,
	}
	client.Write()

	for {

	}
}

func serializeRegisterRequest(ecu_name string, signal_name string, signal_type signals.SIGNAL_TYPE, signal_direction signals.SIGNAL_DIRECTION) []byte {
	builder2 := flatbuffers.NewBuilder(1024)
	ecu2 := builder2.CreateString(ecu_name)
	signal_name2 := builder2.CreateString(signal_name)

	signals.RegisterRequestStart(builder2)
	signals.RegisterRequestAddEcuName(builder2, ecu2)
	signals.RegisterRequestAddSignalName(builder2, signal_name2)
	signals.RegisterRequestAddSignalType(builder2, signal_type)
	signals.RegisterRequestAddSignalDirection(builder2, signal_direction)

	reg_request := signals.RegisterRequestEnd(builder2)
	builder2.Finish(reg_request)
	return builder2.FinishedBytes()
}

func serializeReadRequest(ecu_name string, signal_name string, signal_type signals.SIGNAL_TYPE) []byte {
	builder := flatbuffers.NewBuilder(1024)
	ecu := builder.CreateString(ecu_name)
	sig_name := builder.CreateString(signal_name)

	signals.ReadRequestStart(builder)
	signals.ReadRequestAddEcuName(builder, ecu)
	signals.ReadRequestAddSignalName(builder, sig_name)
	signals.ReadRequestAddSignalType(builder, signals.SIGNAL_TYPEDIGITAL)
	readRequest := signals.ReadRequestEnd(builder)

	signals.RequestStart(builder)
	signals.RequestAddRequestType(builder, signals.RequestTypeReadRequest)
	signals.RequestAddRequest(builder, readRequest)
	request := signals.RequestEnd(builder)
	builder.Finish(request)
	return builder.FinishedBytes()
}

func serializeSetRequest(ecu_name string, signal_name string, signal_type signals.SIGNAL_TYPE, voltage float64, level bool) []byte {
	builder2 := flatbuffers.NewBuilder(1024)
	ecu2 := builder2.CreateString(ecu_name)
	signal_name2 := builder2.CreateString(signal_name)

	switch signal_type {
	case signals.SIGNAL_TYPEDIGITAL:
		signals.DigitalStart(builder2)
		signals.DigitalAddValue(builder2, level)
		dig_sig := signals.DigitalEnd(builder2)

		signals.SetRequestStart(builder2)
		signals.SetRequestAddEcuName(builder2, ecu2)
		signals.SetRequestAddSignalName(builder2, signal_name2)
		signals.SetRequestAddSignalType(builder2, signals.SIGNAL_TYPEDIGITAL)
		signals.SetRequestAddSignalValueType(builder2, signals.SignalValueDigital)
		signals.SetRequestAddSignalValue(builder2, dig_sig)
		setRequest2 := signals.ReadRequestEnd(builder2)

		signals.RequestStart(builder2)
		signals.RequestAddRequestType(builder2, signals.RequestTypeSetRequest)
		signals.RequestAddRequest(builder2, setRequest2)
		setRequest := signals.RequestEnd(builder2)
		builder2.Finish(setRequest)
		return builder2.FinishedBytes()
	case signals.SIGNAL_TYPEANALOG:
		signals.AnalogStart(builder2)
		signals.AnalogAddVoltage(builder2, voltage)
		dig_sig := signals.AnalogEnd(builder2)

		signals.SetRequestStart(builder2)
		signals.SetRequestAddEcuName(builder2, ecu2)
		signals.SetRequestAddSignalName(builder2, signal_name2)
		signals.SetRequestAddSignalType(builder2, signals.SIGNAL_TYPEANALOG)
		signals.SetRequestAddSignalValueType(builder2, signals.SignalValueAnalog)
		signals.SetRequestAddSignalValue(builder2, dig_sig)
		setRequest2 := signals.ReadRequestEnd(builder2)

		signals.RequestStart(builder2)
		signals.RequestAddRequestType(builder2, signals.RequestTypeSetRequest)
		signals.RequestAddRequest(builder2, setRequest2)
		setRequest := signals.RequestEnd(builder2)
		builder2.Finish(setRequest)
		return builder2.FinishedBytes()
	}
	return nil
}
