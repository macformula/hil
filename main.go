package main

import (
	"fmt"
	"sync"
	// "fmt"

	//d "github.com/macformula/hil/dispatcher"
	test "github.com/macformula/hil/clitesting"
)

func main() {
	var wg sync.WaitGroup

	//var g *d.GithubActions
	//var mu sync.Mutex // Mutex to protect access to g

	wg.Add(1)

	//ctx, _ := context.WithCancel(context.Background())
	//g = d.NewGithubActions(nil, 8080)

	go func() {
		defer wg.Done()

		//g.Start(ctx)
	}()

	fmt.Println("passed first wait")

	//mu.Lock()
	//signalGithubActions := g.GetStartSignal(ctx)
	test.Start()
	//mu.Unlock()

	//for {
	//	select {
	//	case <-signalGithubActions:
	//		fmt.Printf("Channel %s received a message!\n", "GithubActions")
	//	}
	//}
	wg.Wait()
}
