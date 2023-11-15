package main

import (
	"context"
	"fmt"
	"sync"
	// "fmt"

	d "github.com/macformula/hil/dispatcher"
)

func main() {
	var wg sync.WaitGroup

	var g *d.GithubActions
	//var mu sync.Mutex // Mutex to protect access to g

	wg.Add(1)

	ctx, _ := context.WithCancel(context.Background())
	g = d.NewGithubActions(nil, 8080)

	go func() {
		defer wg.Done()

		g.Start(ctx)
	}()

	fmt.Println("passed first wait")

	//mu.Lock()
	signal := g.GetStartSignal(ctx)
	//mu.Unlock()

	for {
		<-signal
	}
	wg.Wait()
}
