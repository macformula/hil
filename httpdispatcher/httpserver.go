package httpdispatcher

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	"go.uber.org/zap"
)

type HttpServer struct {
	l                *zap.Logger
	start            chan orchestrator.StartSignal
	results          chan orchestrator.ResultsSignal
	status           chan orchestrator.StatusSignal
	cancelTest       chan orchestrator.CancelTestSignal
	recoverFromFatal chan orchestrator.RecoverFromFatalSignal
	shutdown         chan orchestrator.ShutdownSignal
	sequences        []flow.Sequence
}

func NewHttpServer(l *zap.Logger) *HttpServer {
	return &HttpServer{
		l:                l.Named(_httpLoggerName),
		start:            make(chan orchestrator.StartSignal),
		results:          make(chan orchestrator.ResultsSignal),
		status:           make(chan orchestrator.StatusSignal),
		cancelTest:       make(chan orchestrator.CancelTestSignal),
		recoverFromFatal: make(chan orchestrator.RecoverFromFatalSignal),
		shutdown:         make(chan orchestrator.ShutdownSignal),
	}
}

func (h *HttpServer) Open(ctx context.Context, sequences []flow.Sequence) error {
	h.sequences = sequences

	err := h.setupServer()
	if err != nil {
		return err
	}

	go h.startServer()

	return nil
}

func (h *HttpServer) Close() error {
	err := h.closeServer()
	if err != nil {
		return err
	}

	return nil
}

func (h *HttpServer) Start() <-chan orchestrator.StartSignal {
	return h.start
}

func (h *HttpServer) CancelTest() <-chan orchestrator.CancelTestSignal {
	return h.cancelTest
}

func (h *HttpServer) Shutdown() <-chan orchestrator.ShutdownSignal {
	return h.shutdown
}

func (h *HttpServer) RecoverFromFatal() <-chan orchestrator.RecoverFromFatalSignal {
	return h.recoverFromFatal
}

func (h *HttpServer) Status() chan<- orchestrator.StatusSignal {
	return h.status
}

func (h *HttpServer) Results() chan<- orchestrator.ResultsSignal {
	return h.results
}

// upgrade from http to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Prepares the HTTP server to handle WebSocket upgrade requests
func (h *HttpServer) setupServer() error {
	http.HandleFunc("/dispatcher", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			h.l.Error("Failed to upgrade to websocket", zap.Error(err))
			return
		}
		defer conn.Close()

		// Handle incoming messages in a separate goroutine
		for {

			messageType, message, err := conn.ReadMessage()
			if err != nil {
				h.l.Error("Read error", zap.Error(err))
				break
			}
			// Log the received message
			h.l.Info("Received message", zap.String("message", string(message)))
			prefix := "Sending: "
			msg := append([]byte(prefix), message...)
			// Echo the message back to the client
			if err := conn.WriteMessage(messageType, msg); err != nil {
				h.l.Error("Write error", zap.Error(err))
				return
			}
		}
	})

	return nil
}

func (h *HttpServer) startServer() {
	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

func (h *HttpServer) closeServer() error {
	//TODO
	return nil
}
