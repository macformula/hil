package main

import (
	"sync"

	"github.com/macformula/hil/dispatcher"
)

//var mux sync.Mutex

func main() {
	var wg sync.WaitGroup

	// Increment the wait group counter for each goroutine
	wg.Add(1)

	// Run Dispatcher and ANTHING ELSE as goroutines
	go func() {
		defer wg.Done()
		dispatcher.Dispatcher()
	}()

	wg.Wait()
}
