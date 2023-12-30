package main

import (
	"context"
	"fmt"
	dispatcher "github.com/macformula/hil/dispatcher"
	"go.uber.org/zap"
)

func main() {
	d := dispatcher.NewCliDispatcher(zap.L())
	err := d.Open(context.Background())
	if err != nil {
		return
	}
	start := d.Start()

	fmt.Print(<-start)
	//time.Sleep(5 * time.Second)
	//d.Close()

	//dt := time.Now()
	//<-sig
	//fmt.Println("Current date and time is: ", dt.String())
	//fmt.Print("end")
}
