package main

import (
	"context"
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/cmd/canclienttest/output/CANBMScan"
	"github.com/pkg/errors"
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

	client.Open()

	go send(tx, CANBMScan.NewContactor_Feedback().Frame(), 1, time.Second)
	msg, err := client.Read(context.Background(), CANBMScan.NewContactor_Feedback(), CANBMScan.NewPack_SOC())
	if err != nil {
		logger.Error("client read", zap.Error(err))
	}

	if msg.Frame().ID != CANBMScan.NewContactor_Feedback().Frame().ID {
		logger.Error("client read", zap.Error(errors.New("incorrect CAN frame was read")))
	}

	go send(tx, CANBMScan.NewContactor_States().Frame(), 1, time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	msg, err = client.Read(ctx, CANBMScan.NewPack_Current_Limits())
	if err != nil {
		logger.Error("client read", zap.Error(err))
	}

	if msg != nil {
		logger.Error("client read", zap.Error(errors.New("message should not have been read")))
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

	err = client.Close()
	if err != nil {
		logger.Error("client close", zap.Error(err))
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
