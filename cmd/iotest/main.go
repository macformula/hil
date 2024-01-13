package main

import (
	"fmt"
	"github.com/macformula/hil/iocontrol"
	"github.com/macformula/hil/iocontrol/raspi"
	"time"
)

func main() {
	pin1 := raspi.NewDigitalPin(raspi.Gpio6, iocontrol.Output)

	err := pin1.Open()
	if err != nil {
		fmt.Println(err)
	}

	err = pin1.Write(iocontrol.High)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second)

	err = pin1.Write(iocontrol.Low)
	if err != nil {
		fmt.Println(err)
	}

	pin1.SetDirection(iocontrol.Input)

	fmt.Println(pin1.Read())
}
