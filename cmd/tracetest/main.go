package main

import (
	"context"
	"fmt"
	"github.com/macformula/hil/cmd/canclienttest/output/CANBMScan"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
	"time"

	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/config"
	"go.uber.org/zap"
)

const (
	_configPath = "/opt/macfe/etc/config.yml"
)

func main() {
	config, err := config.NewConfig(_configPath)
	if err != nil {
		panic("error reading config file")
	}

	ctx := context.Background()

	cfg := zap.NewDevelopmentConfig()
	formattedTime := time.Now().Format("2006.01.02_15.04.05")
	fileName := fmt.Sprintf("/opt/macfe/traces/logs/file_%s.log", formattedTime)
	cfg.OutputPaths = []string{fileName}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	conn, err := socketcan.DialContext(context.Background(), "can", "can0")
	if err != nil {
		panic(err)
	}

	client := canlink.NewCANClient(CANBMScan.Messages(), conn)
	tx := socketcan.NewTransmitter(conn)

	tracer := canlink.NewTracer(
		config.CanInterface,
		config.TracerDirectory,
		logger,
		conn,
		canlink.WithBusName(config.BusName),
		canlink.WithTimeout(3*time.Second),
		canlink.WithAscii(client),
		canlink.WithMcap(client),
		canlink.WithCSV(client))

	client.Open()
	err = tracer.Open(ctx)
	if err != nil {
		logger.Error("open tracer", zap.Error(err))
		return
	}

	err = tracer.StartTrace(ctx)
	if err != nil {
		logger.Error("start trace", zap.Error(err))
	}

	// First Test
	go send(tx, CANBMScan.NewContactor_Feedback().Frame(), 50, 10*time.Millisecond)
	msg, err := client.Read(context.Background(), CANBMScan.NewContactor_Feedback(), CANBMScan.NewPack_SOC())
	if err != nil {
		logger.Error("client read", zap.Error(err))
	}

	logger.Info("CAN msg", zap.String("msg.string", msg.String()))
	go send(tx, CANBMScan.NewContactor_Feedback().Frame(), 10, 50*time.Millisecond)
	msg, err = client.Read(context.Background(), CANBMScan.NewContactor_Feedback(), CANBMScan.NewPack_SOC())
	if err != nil {
		logger.Error("client read", zap.Error(err))
	}

	logger.Info("CAN msg", zap.String("msg.string", msg.String()))

	//time.Sleep(5 * time.Second)
	time.Sleep(2 * time.Second)

	err = tracer.StopTrace()
	if err != nil {
		logger.Error("stop trace", zap.Error(err))
	}

	//time.Sleep(2 * time.Second)
	//time.Sleep(10 * time.Second)

	//err = tracer.StartTrace(ctx)
	//if err != nil {
	//	logger.Error("start trace", zap.Error(err))
	//}
	//
	////time.Sleep(2 * time.Second)
	//time.Sleep(10 * time.Second)
	//
	//err = tracer.StopTrace()
	//if err != nil {
	//	logger.Error("stop trace", zap.Error(err))
	//}

	err = tracer.Close()
	if err != nil {
		logger.Error("close tracer", zap.Error(err))
	}

	client.Close()

	if tracer.Error() != nil {
		logger.Error("tracer error", zap.Error(tracer.Error()))
	}
}

func send(tx *socketcan.Transmitter, frame can.Frame, n int, delay time.Duration) {
	for i := 0; i < n; i++ {
		time.Sleep(delay)
		if err := tx.TransmitFrame(context.Background(), frame); err != nil {
			panic(err)
		}
	}
}
