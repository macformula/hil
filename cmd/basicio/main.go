package main

import (
	"fmt"
	"io"
	"net"
	"time"

	signals "github.com/macformula/hil/cmd/basicio/signals"

	flatbuffers "github.com/google/flatbuffers/go"
	"go.uber.org/zap"
)

const (
	_port = 8080
)

type Firmware struct {
	l *zap.Logger
}

func (c *Firmware) Write(conn net.Conn) {
	reg_req := serializeRegisterRequest("lv", "raspi_en", signals.SIGNAL_TYPEDIGITAL, signals.SIGNAL_DIRECTIONINPUT)
	set_req1 := serializeSetRequest("lv", "raspi_en", signals.SIGNAL_TYPEDIGITAL, 0.0, false)
	read_req1 := serializeReadRequest("lv", "raspi_en", signals.SIGNAL_TYPEDIGITAL, signals.SIGNAL_DIRECTIONINPUT)
	set_req2 := serializeSetRequest("lv", "raspi_en", signals.SIGNAL_TYPEDIGITAL, 0.0, true)
	read_req2 := serializeReadRequest("lv", "raspi_en", signals.SIGNAL_TYPEDIGITAL, signals.SIGNAL_DIRECTIONINPUT)

	// Send data to the server
	_, err := conn.Write(reg_req)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	time.Sleep(2 * time.Second)

	// Set pin to low
	_, err = conn.Write(set_req1)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	time.Sleep(2 * time.Second)

	// Read pin
	_, err = conn.Write(read_req1)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	time.Sleep(2 * time.Second)

	// Set pin to high
	_, err = conn.Write(set_req2)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	time.Sleep(2 * time.Second)

	// Read pin
	_, err = conn.Write(read_req2)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	time.Sleep(2 * time.Second)
}

func (c *Firmware) Listen(conn net.Conn) {
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

	// Connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", _port))
	if err != nil {
		logger.Info(fmt.Sprintf("Error connecting to server: %e", err))
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	firmware := Firmware{
		l: logger,
	}
	go firmware.Write(conn)
	go firmware.Listen(conn)

	time.Sleep(20 * time.Second)
}

func serializeRegisterRequest(ecu_name string, signal_name string, signalType signals.SIGNAL_TYPE, signal_direction signals.SIGNAL_DIRECTION) []byte {
	builder := flatbuffers.NewBuilder(1024)
	ecu2 := builder.CreateString(ecu_name)
	signal_name2 := builder.CreateString(signal_name)

	signals.RegisterRequestStart(builder)
	signals.RegisterRequestAddEcuName(builder, ecu2)
	signals.RegisterRequestAddSignalName(builder, signal_name2)
	signals.RegisterRequestAddSignalType(builder, signalType)
	signals.RegisterRequestAddSignalDirection(builder, signal_direction)

	registerRequest := signals.RegisterRequestEnd(builder)

	signals.RequestStart(builder)
	signals.RequestAddRequestType(builder, signals.RequestTypeRegisterRequest)
	signals.RequestAddRequest(builder, registerRequest)
	request := signals.RequestEnd(builder)
	builder.Finish(request)
	return builder.FinishedBytes()
}

func serializeReadRequest(ecu_name string, signal_name string, signalType signals.SIGNAL_TYPE, sigDirection signals.SIGNAL_DIRECTION) []byte {
	builder := flatbuffers.NewBuilder(1024)
	ecu := builder.CreateString(ecu_name)
	sig_name := builder.CreateString(signal_name)

	signals.ReadRequestStart(builder)
	signals.ReadRequestAddEcuName(builder, ecu)
	signals.ReadRequestAddSignalName(builder, sig_name)
	signals.ReadRequestAddSignalType(builder, signalType)
	signals.ReadRequestAddSignalDirection(builder, sigDirection)
	readRequest := signals.ReadRequestEnd(builder)

	signals.RequestStart(builder)
	signals.RequestAddRequestType(builder, signals.RequestTypeReadRequest)
	signals.RequestAddRequest(builder, readRequest)
	request := signals.RequestEnd(builder)
	builder.Finish(request)
	return builder.FinishedBytes()
}

func serializeSetRequest(ecu_name string, signal_name string, signalType signals.SIGNAL_TYPE, voltage float64, level bool) []byte {
	builder := flatbuffers.NewBuilder(1024)
	ecu2 := builder.CreateString(ecu_name)
	signal_name2 := builder.CreateString(signal_name)

	switch signalType {
	case signals.SIGNAL_TYPEDIGITAL:
		signals.DigitalStart(builder)
		signals.DigitalAddValue(builder, level)
		dig_sig := signals.DigitalEnd(builder)

		signals.SetRequestStart(builder)
		signals.SetRequestAddEcuName(builder, ecu2)
		signals.SetRequestAddSignalName(builder, signal_name2)
		signals.SetRequestAddSignalType(builder, signals.SIGNAL_TYPEDIGITAL)
		signals.SetRequestAddSignalValueType(builder, signals.SignalValueDigital)
		signals.SetRequestAddSignalValue(builder, dig_sig)
		setRequest2 := signals.ReadRequestEnd(builder)

		signals.RequestStart(builder)
		signals.RequestAddRequestType(builder, signals.RequestTypeSetRequest)
		signals.RequestAddRequest(builder, setRequest2)
		setRequest := signals.RequestEnd(builder)
		builder.Finish(setRequest)
		return builder.FinishedBytes()
	case signals.SIGNAL_TYPEANALOG:
		signals.AnalogStart(builder)
		signals.AnalogAddVoltage(builder, voltage)
		dig_sig := signals.AnalogEnd(builder)

		signals.SetRequestStart(builder)
		signals.SetRequestAddEcuName(builder, ecu2)
		signals.SetRequestAddSignalName(builder, signal_name2)
		signals.SetRequestAddSignalType(builder, signals.SIGNAL_TYPEANALOG)
		signals.SetRequestAddSignalValueType(builder, signals.SignalValueAnalog)
		signals.SetRequestAddSignalValue(builder, dig_sig)
		setRequest2 := signals.ReadRequestEnd(builder)

		signals.RequestStart(builder)
		signals.RequestAddRequestType(builder, signals.RequestTypeSetRequest)
		signals.RequestAddRequest(builder, setRequest2)
		setRequest := signals.RequestEnd(builder)
		builder.Finish(setRequest)
		return builder.FinishedBytes()
	}
	return nil
}

func deserializeReadResponse(unionTable *flatbuffers.Table) (bool, string, bool, float64) {
	unionResponse := new(signals.ReadResponse)
	unionResponse.Init(unionTable.Bytes, unionTable.Pos)

	ok := unionResponse.Ok()
	errorString := string(unionResponse.Error())

	unionTable = new(flatbuffers.Table)
	if unionResponse.SignalValue(unionTable) {
		switch unionResponse.SignalValueType() {
		case signals.SignalValueDigital:
			unionSignalValue := new(signals.Digital)
			unionSignalValue.Init(unionTable.Bytes, unionTable.Pos)

			return ok, errorString, unionSignalValue.Value(), 0.0

		case signals.SignalValueAnalog:
			unionSignalValue := new(signals.Analog)
			unionSignalValue.Init(unionTable.Bytes, unionTable.Pos)

			return ok, errorString, false, unionSignalValue.Voltage()
		}
	}

	return false, "", false, 0.0
}
