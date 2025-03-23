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

	builder := flatbuffers.NewBuilder(1024)
	ecu := builder.CreateString("lv_controller")
	signal_name := builder.CreateString("raspi_en")

	signals.ReadRequestStart(builder)
	signals.ReadRequestAddEcuName(builder, ecu)
	signals.ReadRequestAddSignalName(builder, signal_name)
	signals.ReadRequestAddSignalType(builder, signals.SIGNAL_TYPEDIGITAL)
	readRequest := signals.ReadRequestEnd(builder)

	signals.RequestStart(builder)
	signals.RequestAddRequestType(builder, signals.RequestTypeReadRequest)
	signals.RequestAddRequest(builder, readRequest)
	request := signals.RequestEnd(builder)
	builder.Finish(request)
	buf := builder.FinishedBytes()

	builder2 := flatbuffers.NewBuilder(1024)
	ecu2 := builder2.CreateString("tms")
	signal_name2 := builder2.CreateString("signal 2")

	signals.DigitalStart(builder2)
	signals.DigitalAddValue(builder2, true)
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
	buf2 := builder2.FinishedBytes()

	// Send data to the server
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	time.Sleep(2 * time.Second)

	_, err = conn.Write(buf2)
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
