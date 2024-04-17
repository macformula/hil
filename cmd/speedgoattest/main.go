package main

import (
	"time"

	"github.com/macformula/hil/iocontrol/speedgoat"
	"go.uber.org/zap"
)

const _speedgoatAddr = "192.168.7.5:8001"

func main() {
	cfg := zap.NewDevelopmentConfig()
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	controller := speedgoat.NewController(logger, _speedgoatAddr)

	err = controller.Open()
	if err != nil {
		logger.Error("controller open", zap.Error(err))
		return
	}

	pin := speedgoat.NewDigitalPin(8) // First digital output pin (idx 8-15)
	controller.SetDigital(pin, true)

	pin2 := speedgoat.NewDigitalPin(15) // Last digital output pin
	controller.SetDigital(pin2, true)

	pin3 := speedgoat.NewAnalogPin(8) // First analog output pin (idx 8-11)
	controller.WriteVoltage(pin3, 2.5)

	time.Sleep(time.Second * 2)

	// Verify that these change on the Speedgoat via Simulink or LED
	// If they don't change (or take a while), verify that the Simulink TCP server has a sufficiently small sample time
	controller.SetDigital(pin2, false)
	controller.WriteVoltage(pin3, 3)

	time.Sleep(time.Second * 2)

	err = controller.Close()
	if err != nil {
		logger.Error("controller close", zap.Error(err))
		return
	}
}
