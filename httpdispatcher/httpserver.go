package httpdispatcher

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"

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
	sequences        map[int]flow.Sequence
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
	// Create a map of Sequences with an int ID
	seqMap := make(map[int]flow.Sequence)
	for i, seq := range sequences {
		seqMap[i] = seq
	}
	h.sequences = seqMap

	h.setupServer()
	// if err != nil {
	// 	return err
	// }

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
	CheckOrigin:     func(r *http.Request) bool { return true },
	// replace return true with the following after testing:
	// allowedOrigins := []string{"https://dev.macformularacing.com"}
	// origin := r.Header.Get("Origin")
	//     for _, o := range allowedOrigins {
	//         if origin == o {
	//             return true
	//         }
	//     }
	//     return false
	// },
}

func (h *HttpServer) createWS(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.l.Error("Failed to upgrade to websocket", zap.Error(err))
		return nil
	}
	return conn
}

func (h *HttpServer) setupServer() {
	http.HandleFunc("/status", h.serveStatus)
	http.HandleFunc("/sequences", h.serveSequences)
	http.HandleFunc("/start", h.serveStart) // results will be sent here
	// http.HandleFunc("/cancel", serveStatus)
	// http.HandleFunc("/recover", serveStatus)
}

func (h *HttpServer) serveStatus(w http.ResponseWriter, r *http.Request) {
	conn := h.createWS(w, r)
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			h.l.Error("Read error", zap.Error(err))
			break
		}
		// Log the received message
		h.l.Info("Received message", zap.String("message", string(message)))

		statusJSON, err := json.Marshal(h.Status())
		if err := conn.WriteMessage(messageType, statusJSON); err != nil {
			h.l.Error("Write error", zap.Error(err))
			return
		}
	}
}

func (h *HttpServer) serveStart(w http.ResponseWriter, r *http.Request) {
	conn := h.createWS(w, r)
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			h.l.Error("Read error", zap.Error(err))
			break
		}

		messageInt, err := strconv.Atoi(string(message))
		if err != nil {
			h.l.Error("The message is not an integer.", zap.Error(err))
			break
		}

		// Log the received message
		h.l.Info("Received message", zap.String("message", string(message)))
		var httpCode string
		if sequence, ok := h.sequences[messageInt]; ok {
			// key exists in sequences map
			h.start <- orchestrator.StartSignal{
				TestId:   uuid.New(),
				Seq:      sequence,
				Metadata: nil,
			}
			httpCode = "204"
		} else {
			httpCode = "400"
		}

		// Echo the message back to the client
		if err := conn.WriteMessage(messageType, []byte(httpCode)); err != nil {
			h.l.Error("Write error", zap.Error(err))
			return
		}
	}
}

func (h *HttpServer) serveSequences(w http.ResponseWriter, r *http.Request) {
	conn := h.createWS(w, r)
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			h.l.Error("Read error", zap.Error(err))
			break
		}
		// Log the received message
		h.l.Info("Received message", zap.String("message", string(message)))

		// send sequences
		sequencesJSON, err := json.Marshal(h.sequences)
		if err := conn.WriteMessage(messageType, sequencesJSON); err != nil {
			h.l.Error("Write error", zap.Error(err))
			return
		}
	}
}

// func (h *HttpServer) serveCancel(w http.ResponseWriter, r *http.Request) {
// 	conn := h.createWS(w, r)

// }

// func (h *HttpServer) serverRecover(w http.ResponseWriter, r *http.Request) {
// 	conn := h.createWS(w, r)

// }

// func (h *HttpServer) setupServer() error {
// 	err := h.setupDispatcherStatusEndpoint()
// 	if err != nil {
// 		return err
// 	}

// 	err = h.setupSequencesEndpoint()
// 	if err != nil {
// 		return err
// 	}

// 	err = h.setupStartTestEndpoint()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// Prepares the HTTP server to handle WebSocket upgrade requests
// func (h *HttpServer) setupDispatcherStatusEndpoint() error {
// 	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
// 		conn := h.createWS(w, r)
// 		defer conn.Close()

// 		// Handle incoming messages in a separate goroutine
// 		for {

// 			messageType, message, err := conn.ReadMessage()
// 			if err != nil {
// 				h.l.Error("Read error", zap.Error(err))
// 				break
// 			}
// 			// Log the received message
// 			h.l.Info("Received message", zap.String("message", string(message)))
// 			prefix := "Sending: "
// 			msg := append([]byte(prefix), message...)
// 			// Echo the message back to the client
// 			if err := conn.WriteMessage(messageType, msg); err != nil {
// 				h.l.Error("Write error", zap.Error(err))
// 				return
// 			}
// 		}
// 	})

// 	return nil
// }

// func (h *HttpServer) setupSequencesEndpoint() error {
// 	http.HandleFunc("/sequences", func(w http.ResponseWriter, r *http.Request) {
// 		conn := h.createWS(w, r)
// 		defer conn.Close()

// 		// Handle incoming messages in a separate goroutine
// 		for {

// 			messageType, message, err := conn.ReadMessage()
// 			if err != nil {
// 				h.l.Error("Read error", zap.Error(err))
// 				break
// 			}
// 			// Log the received message
// 			h.l.Info("Received message", zap.String("message", string(message)))

// 			jsonData, err := json.Marshal(h.sequences)

// 			// Echo the message back to the client
// 			if err := conn.WriteMessage(messageType, jsonData); err != nil {
// 				h.l.Error("Write error", zap.Error(err))
// 				return
// 			}
// 		}
// 	})

// 	return nil
// }

// func (h *HttpServer) setupStartTestEndpoint() error {
// 	http.HandleFunc("/starttest", func(w http.ResponseWriter, r *http.Request) {
// 		conn, err := upgrader.Upgrade(w, r, nil)
// 		if err != nil {
// 			h.l.Error("Failed to upgrade to websocket", zap.Error(err))
// 			return
// 		}
// 		defer conn.Close()

// 		// Handle incoming messages in a separate goroutine
// 		for {

// 			messageType, message, err := conn.ReadMessage()
// 			if err != nil {
// 				h.l.Error("Read error", zap.Error(err))
// 				break
// 			}

// 			messageInt, err := strconv.Atoi(string(message))
// 			if err != nil {
// 				h.l.Error("The message is not an integer.", zap.Error(err))
// 				break
// 			}

// 			// Log the received message
// 			h.l.Info("Received message", zap.String("message", string(message)))
// 			var httpCode string
// 			if sequence, ok := h.sequences[messageInt]; ok {
// 				// key exists in sequences map
// 				h.start <- orchestrator.StartSignal{
// 					TestId:   uuid.New(),
// 					Seq:      sequence,
// 					Metadata: nil,
// 				}
// 				httpCode = "204"
// 			} else {
// 				httpCode = "400"
// 			}

// 			// Echo the message back to the client
// 			if err := conn.WriteMessage(messageType, []byte(httpCode)); err != nil {
// 				h.l.Error("Write error", zap.Error(err))
// 				return
// 			}
// 		}
// 	})

// 	return nil
// }

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
