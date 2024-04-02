package main

import (
	"go.uber.org/zap"
	"time"
)

func main() {
	controller := NewController(zap.L(), "192.168.7.5:8001")

	err := controller.Open()
	if err != nil {
		panic(err)
	}

	pin := NewDigitalPin(0)
	controller.SetDigital(pin, true)
	pin2 := NewDigitalPin(15)
	controller.SetDigital(pin2, true)
	pin3 := NewAnalogPin(8)
	controller.WriteVoltage(pin3, 1000)
	println("set!")
	time.Sleep(time.Millisecond * 2000)
	println("setting false")
	controller.SetDigital(pin2, false)
	println("set false")
	//controller.WriteVoltage(pin3, 3)
	time.Sleep(time.Second * 2)
}
