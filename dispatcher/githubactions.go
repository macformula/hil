package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type GithubActions struct {
	l         *zap.Logger
	port      int
	srv       *http.Server
	startChan chan struct{}
	stop      bool
}

func NewGithubActions(l *zap.Logger, port int) *GithubActions {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.DefaultServeMux,
	}

	return &GithubActions{
		l:         l,
		port:      port,
		srv:       server,
		startChan: make(chan struct{}),
		stop:      true,
	}
}

func (g *GithubActions) Start(ctx context.Context) error {
	if g.stop {
		g.stop = false

		http.HandleFunc("/trigger", func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-ctx.Done():
				return
			default:
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Error reading request body", http.StatusBadRequest)
					return
				}

				fmt.Printf("Received a request with message: %s\n", string(body))

				g.startChan <- struct{}{}
				fmt.Println("unblocked on channel")
				// Set status code
				w.WriteHeader(http.StatusAccepted)
				//w.Write([]byte(message))

				// https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(NewFakeTestData())
				fmt.Println("MESSAGE SENT BACK")
			}
		})

		go func() {
			// fmt.Printf("Listening on :%d...\n", d.port)
			if err := g.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Println("Error:", err)
			}
		}()
		fmt.Println("got here")
		// block until ctx expired or stop called
		select {
		case <-ctx.Done():
			// Context is canceled, proceed to shutdown the server
			g.Stop(ctx)
			fmt.Println("closing")
		}
	}

	return nil
}

func (g *GithubActions) GetStartSignal(ctx context.Context) <-chan struct{} {
	return g.startChan
}

func (g *GithubActions) SetResults(ctx context.Context) error {
	return nil
}

func (g *GithubActions) Stop(ctx context.Context) error {
	if !g.stop {
		fmt.Println("Shutting down...")
		g.stop = true
		return g.srv.Shutdown(ctx)
	}
	return nil
}
