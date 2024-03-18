package httpdispatcher

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/event"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	"go.uber.org/zap"
)

type Message struct {
	Task      string `json:"task"`
	Parameter string `json:"parameter"`
}

// message task values
const (
	StartTest        = "start"
	CancelTest       = "cancel"
	RecoverFromFatal = "recover"
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
	statusFeed       event.Feed
	resultsFeed      event.Feed
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
		statusFeed:       event.Feed{},
		resultsFeed:      event.Feed{},
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

	go h.monitorDispatcher(ctx)

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

func (h *HttpServer) SubscribeToFeeds(progressChan chan orchestrator.StatusSignal, resultsChan chan orchestrator.ResultsSignal) (progressSub event.Subscription, resultSub event.Subscription) {
	return h.statusFeed.Subscribe(progressChan), h.resultsFeed.Subscribe(resultsChan) // TODO make results/status feed - update in monitordispatcher
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

// Websocket endpoints
func (h *HttpServer) setupServer() {
	http.HandleFunc("/test", h.serveTest) // start/sequences - cancel - recover
	// http.HandleFunc("/status", h.serveStatus)
}

func (h *HttpServer) serveTest(w http.ResponseWriter, r *http.Request) {
	conn := h.createWS(w, r)
	client := NewClient(conn)
	progressSub, resultsSub := h.SubscribeToFeeds(client.status, client.results)

	go client.updateTests()

	defer client.conn.Close()
	defer progressSub.Unsubscribe()
	defer resultsSub.Unsubscribe()

	for {
		msg := h.readWS(conn)

		switch msg.Task {
		case StartTest:
			h.startClientTest(client, msg.Parameter)
		case CancelTest:
			h.cancelClientTest(client, msg.Parameter)
		case RecoverFromFatal:
			h.recoverClientFromFatal(client)
		}
	}
}

func (h *HttpServer) startClientTest(client *Client, parameter string) {
	var messageInt int
	var err error

	// If no parameter is provided, send the list of sequences to the client.
	if parameter == "" {
		sequencesJSON, _ := json.Marshal(h.sequences)
		if err := client.conn.WriteMessage(websocket.TextMessage, sequencesJSON); err != nil {
			h.l.Error("Write error", zap.Error(err))
		}
		return
	} else {
		// If a parameter is provided, convert it to an integer.
		messageInt, err = strconv.Atoi(string(parameter))
		if err != nil {
			h.l.Error("The message is not an integer.", zap.Error(err))
			return
		}
	}

	if sequence, ok := h.sequences[messageInt]; ok {
		// Send start signal
		newTestID := uuid.New()
		h.start <- orchestrator.StartSignal{
			TestId:   newTestID,
			Seq:      sequence,
			Metadata: nil,
		}
		// Add test to client test list
		queuePosition := (<-client.status).QueueLength
		client.addTest(queuePosition, newTestID)

		// Send queue position and new test ID
		testData, _ := json.Marshal(strconv.Itoa(queuePosition) + ", " + newTestID.String())
		client.conn.WriteMessage(websocket.TextMessage, testData)
	} else {
		badRequest, _ := json.Marshal(http.StatusBadRequest)
		client.conn.WriteMessage(websocket.TextMessage, badRequest)
	}
}

func (h *HttpServer) cancelClientTest(client *Client, parameter string) {
	var testIndex int
	var err error

	// Send client cancellable tests
	if parameter == "" {
		testQueueJSON, _ := json.Marshal(client.testQueue)
		if err := client.conn.WriteMessage(websocket.TextMessage, testQueueJSON); err != nil {
			h.l.Error("Write error", zap.Error(err))
		}
		return
	} else {
		testIndex, err = strconv.Atoi(string(parameter))
		if err != nil {
			h.l.Error("The message is not an integer.", zap.Error(err))
			return
		}
	}

	var httpCode string
	if testIndex >= 0 && testIndex < len(client.testQueue) {
		h.cancelTest <- orchestrator.CancelTestSignal{
			TestId: client.testQueue[testIndex].UUID,
		}
		client.removeTest(testIndex)
		httpCode = http.StatusText(http.StatusOK)
	} else {
		httpCode = http.StatusText(http.StatusBadRequest)
	}
	httpCodeJSON, _ := json.Marshal(httpCode)
	client.conn.WriteMessage(websocket.TextMessage, httpCodeJSON)
}

func (h *HttpServer) recoverClientFromFatal(client *Client) {
	h.recoverFromFatal <- orchestrator.RecoverFromFatalSignal{}
	client.conn.WriteMessage(websocket.TextMessage, []byte{http.StatusOK})
}

func (h *HttpServer) readWS(conn *websocket.Conn) *Message {
	_, message, err := conn.ReadMessage()
	if err != nil {
		h.l.Error("Read error", zap.Error(err))
		return nil
	}
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		h.l.Error("JSON Unmarshal error", zap.Error(err))
		return nil
	}
	h.l.Info("Extracted values", zap.String("task", msg.Task), zap.String("parameter", msg.Parameter))
	return &msg
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

func (h *HttpServer) monitorDispatcher(ctx context.Context) {
	for {
		select {
		case stats := <-h.status:
			h.l.Info("status signal received")
		    h.statusFeed.Send(stats)

		case res := <-h.results:
			h.l.Info("results signal received")
			h.resultsFeed.Send(res)
		case <-ctx.Done():
			h.l.Info("context done signal received")

			return
		}
	}
}
