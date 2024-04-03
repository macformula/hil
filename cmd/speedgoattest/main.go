package main

import (
	"github.com/macformula/hil/iocontrol/speedgoat"
	"go.uber.org/zap"
	"time"
)

func main() {
	controller := speedgoat.NewController(zap.L(), "192.168.7.5:8001")

	err := controller.Open()
	if err != nil {
		panic(err)
	}

	pin := speedgoat.NewDigitalPin(8)
	controller.SetDigital(pin, true)
	pin2 := speedgoat.NewDigitalPin(15)
	controller.SetDigital(pin2, true)
	pin3 := speedgoat.NewAnalogPin(8)
	controller.WriteVoltage(pin3, 1000)

	time.Sleep(time.Second * 2)
	
	controller.SetDigital(pin2, false)
	controller.WriteVoltage(pin3, 3)

	time.Sleep(time.Second * 2)
}
