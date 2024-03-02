package main

import (
	"context"
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/cmd/canclienttest/output/CANBMScan"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
	"time"
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	conn, err := socketcan.DialContext(context.Background(), "can", "can0")
	if err != nil {
		panic(err)
	}

	client := canlink.NewCANClient(CANBMScan.Messages(), conn)
	tx := socketcan.NewTransmitter(conn)

	go send(tx, CANBMScan.NewContactor_Feedback().Frame(), 1, time.Second)
	_, err = client.Read(context.Background())
	if err != nil {
		logger.Error("client read", zap.Error(err))
	}

	err = client.StartTracking()
	if err != nil {
		logger.Error("start tracking", zap.Error(err))
	}

	go send(tx, CANBMScan.NewContactor_States().Frame(), 4, time.Millisecond*10)
	go send(tx, CANBMScan.NewPack_SOC().Frame(), 6, time.Millisecond*10)
	time.Sleep(time.Second)

	data, err := client.StopTracking()
	if err != nil {
		logger.Error("stop tracking", zap.Error(err))
	}

	logger.Info("tracker", zap.Uint32("contactor_states frames", data[CANBMScan.NewContactor_States().Frame().ID]))
	logger.Info("tracker", zap.Uint32("pack_soc frames", data[CANBMScan.NewPack_SOC().Frame().ID]))
}

func send(tx *socketcan.Transmitter, frame can.Frame, n int, delay time.Duration) {
	for i := 0; i < n; i++ {
		time.Sleep(delay)
		if err := tx.TransmitFrame(context.Background(), frame); err != nil {
			panic(err)
		}
	}
}
