package main

import (
	"sync"
	"context"
	"fmt"
	"time"

	"github.com/macformula/hil/dispatcher"
)

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Increment the wait group counter for each goroutine
	wg.Add(1)

	// Run Dispatcher and ANTHING ELSE as goroutines
	d := dispatcher.NewDispatcher(nil, 8080)
	go func() {
		defer wg.Done()
		
		if err := d.Open(ctx); err != nil {
			fmt.Println("Error opening server:", err)
			return
		}
	}()
	time.Sleep(1 * time.Second)
	cancel()

	time.Sleep(1 * time.Second)

	if err := d.Close(ctx); err != nil {
		fmt.Println("Error closing server:", err)
	}

	wg.Wait()
}
