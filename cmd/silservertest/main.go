package main

import (
	"context"
	"fmt"
	"net"

	signals "github.com/macformula/hil/cmd/silservertest/signals"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/macformula/hil/iocontrol/sil"
	"go.uber.org/zap"
)

type Client struct {
	addr string
	l *zap.Logger
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

	// Send data to the server
	_, err = conn.Write(buf)
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
		l: logger,
	}
	client.Write()

	for {

	}
}