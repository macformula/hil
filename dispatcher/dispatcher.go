package dispatcher

import (
	"fmt"
	"io"
	"net/http"
	// "time"
	"context"
	// "os"
	// "os/exec"
	// "runtime"

	"github.com/macformula/hil/cli"
	"github.com/macformula/hil/orchestrator"
	"go.uber.org/zap"
)

type Dispatcher struct {
	l *zap.Logger
	port int
	srv *http.Server
	stopChan chan struct{}
	stop bool
}

func NewDispatcher(l *zap.Logger, port int) *Dispatcher {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.DefaultServeMux,
	}

	return &Dispatcher{
		l: l,
		port: port,
		srv: server,
		// stopChan: make(chan struct{}),
		stop: true,
	}
}

func (d *Dispatcher) Open(ctx context.Context) error {
	// can only set this up once
	if d.stop {
		d.stop = false
		d.stopChan = make(chan struct{})
		go d.setupCLI() // start cli

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-ctx.Done():
				return
			case <-d.stopChan:
				return
			default:
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Error reading request body", http.StatusBadRequest)
					return
				}

				message := string(body) + " returning..."

				fmt.Printf("Received a request with message: %s\n", message)

				d.Trigger(ctx)

				// Set status code
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte(message))
			}
		})

		go func() {
			// fmt.Printf("Listening on :%d...\n", d.port)
			if err := d.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Println("Error:", err)
			}
		}()

		// block until ctx expired or stop called
		select {
		case <-ctx.Done():
			// Context is canceled, proceed to shutdown the server
			d.Close(ctx)
		case <-d.stopChan:
			// Stop signal received, immediately return
		}
	}

	return nil
}

func (d *Dispatcher) setupCLI() {
	items := make(map[string]func())
	items["Pizza"] = func() { 
		go d.setupCLI()
		d.Open(context.TODO())
	}
	items["Sushi"] = func() { 
		go d.setupCLI()
		fmt.Println("At Sushi")
	}
	items["Burger"] = func() { 
		go d.setupCLI()
		fmt.Println("At Burger") 
	}

	cli.Start(items)
}

func (d *Dispatcher) Trigger(ctx context.Context) error {
	o := orchestrator.NewOrchestrator(nil)
	o.StartTests(ctx)
	return nil
}

// shut down the HTTP server gracefully
func (d *Dispatcher) Close(ctx context.Context) error {
	if !d.stop {
		fmt.Println("Shutting down...")
		d.stop = true
		close(d.stopChan)
		return d.srv.Shutdown(ctx)
	}
	return nil
}
