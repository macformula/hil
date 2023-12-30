package main

import (
	"context"
	dispatcher "github.com/macformula/hil/dispatcher"
	"time"
)

func main() {
	d := dispatcher.NewCliDispatcher()
	err := d.Open(context.Background())
	if err != nil {
		return
	}
	d.Start()

	time.Sleep(5 * time.Second)
	//d.Close()

	//dt := time.Now()
	//<-sig
	//fmt.Println("Current date and time is: ", dt.String())
	//fmt.Print("end")
}
