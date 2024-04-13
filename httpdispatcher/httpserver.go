package httpdispatcher

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

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

type StatusMessage struct {
	Message string
	Code    string
}

type TestQueueItem struct {
	SequenceName 	string
	UUID 			uuid.UUID
	client 			*Client
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
	// testQueueUpdated chan bool
	sequences        map[int]flow.Sequence
	testQueue		 []TestQueueItem
	statusFeed       event.Feed
	resultsFeed      event.Feed
	testQueueUpdateFeed  event.Feed
}

func NewHttpServer(l *zap.Logger) *HttpServer {
	return &HttpServer{
		l:                		l.Named(_httpLoggerName),
		start:            		make(chan orchestrator.StartSignal),
		results:          		make(chan orchestrator.ResultsSignal),
		status:           		make(chan orchestrator.StatusSignal),
		cancelTest:       		make(chan orchestrator.CancelTestSignal),
		recoverFromFatal: 		make(chan orchestrator.RecoverFromFatalSignal),
		shutdown:         		make(chan orchestrator.ShutdownSignal),
		// testQueueUpdated:		make(chan bool),
		testQueue:		  		make([]TestQueueItem, 0),
		statusFeed:       		event.Feed{},
		resultsFeed:      		event.Feed{}, // TODO make results/status feed - update in monitordispatcher
		testQueueUpdateFeed:  	event.Feed{},
	}
}

func (h *HttpServer) Open(ctx context.Context, sequences []flow.Sequence) error {
	// Create a map of Sequences with an int ID
	seqMap := make(map[int]flow.Sequence)
	for i, seq := range sequences {
		seqMap[i] = seq
	}
	h.sequences = seqMap

	go h.monitorDispatcher(ctx)

	go h.StartServer()

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

func (h *HttpServer) StartServer() {
	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/test", h.serveTest) // start/sequences - cancel - recover
	mux.HandleFunc("/status", h.serveStatus)
	mux.HandleFunc("/results", h.serveResults)
	mux.HandleFunc("/queue", h.serveQueue)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

func (h *HttpServer) createWS(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.l.Error("Failed to upgrade to websocket", zap.Error(err))
		return nil
	}
	return conn
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (h *HttpServer) serveTest(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	conn := h.createWS(w, r)
	client := NewClient(conn)

	defer client.conn.Close()
	for {
		msg, err := h.readWS(conn)
		status := StatusMessage{Code: "200", Message: ""}
		if err != nil {
			status.Code = "400"
		}

		switch msg.Task {
		case StartTest:
			h.startClientTest(client, msg.Parameter)
			status.Message = "Client Test Started"
		case CancelTest:
			h.cancelClientTest(client, msg.Parameter)
			status.Message = "Client Test Cancelled"
		case RecoverFromFatal:
			h.recoverClientFromFatal(client)
			status.Message = "Recovered From Fatal"
		default:
			h.l.Info("serveTest Invalid Message Received")
			status.Message = "Invalid Message Received"
		}

		//err = conn.WriteMessage(status)
		statusJSON, _ := json.Marshal(status)
		conn.WriteMessage(websocket.TextMessage, statusJSON)
		if err != nil {
			h.l.Error(errors.Wrap(err, "couldn't send back websocket message").Error())
		}
	}
}

func (h *HttpServer) serveStatus(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	conn := h.createWS(w, r)
	client := NewClient(conn)
	progressSub := h.statusFeed.Subscribe(client.status)

	defer client.conn.Close()
	defer progressSub.Unsubscribe()

	for {
		select {
		case currentStatus := <-client.status:
			currentStatusJSON, _ := json.Marshal(currentStatus)
			conn.WriteMessage(websocket.TextMessage, currentStatusJSON)
		}
	}
}

func (h *HttpServer) serveResults(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	conn := h.createWS(w, r)
	client := NewClient(conn)
	resultSub := h.resultsFeed.Subscribe(client.results)

	defer client.conn.Close()
	defer resultSub.Unsubscribe()

	for {
		select {
		case result := <-client.results:
			resultJSON, _ := json.Marshal(result)
			conn.WriteMessage(websocket.TextMessage, resultJSON)
		}
	}
}

func (h *HttpServer) serveQueue(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	conn := h.createWS(w, r)
	client := NewClient(conn)
	updateSub := h.testQueueUpdateFeed.Subscribe(client.testQueueUpdated)

	defer client.conn.Close()
	defer updateSub.Unsubscribe()

	for {
		select {
		case <- client.testQueueUpdated:
			testQueueJSON, _ := json.Marshal(h.testQueue)
			conn.WriteMessage(websocket.TextMessage, testQueueJSON)
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
		messageInt, err = strconv.Atoi(parameter)
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
		// Add test to test list
		queuePosition := (<-client.status).QueueLength		
		h.addTestToQueue(newTestID, messageInt, client)

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
		h.removeTestFromQueue(client.testQueue[testIndex].UUID)
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

func (h *HttpServer) readWS(conn *websocket.Conn) (*Message, error) {
	messageType, message, err := conn.ReadMessage()
	h.l.Info("Inside readWS")
	if err != nil {
		h.l.Error("Read error", zap.Error(err), zap.Any("message", message), zap.Any("messageType", messageType))
		return &Message{
			Task:      "",
			Parameter: "",
		}, err
	}
	var msg Message
	if err = json.Unmarshal(message, &msg); err != nil {
		h.l.Error("JSON Unmarshal error", zap.Error(err), zap.Any("message", string(message)), zap.Any("messageType", messageType))
		return &Message{
			Task:      "",
			Parameter: "",
		}, err
	}
	h.l.Info("Extracted values", zap.String("task", msg.Task), zap.String("parameter", msg.Parameter))
	return &msg, nil
}

func (h *HttpServer) addTestToQueue(testID uuid.UUID, testIndex int, client *Client) {
	newItem := TestQueueItem{
		SequenceName: h.sequences[testIndex].Name,
		UUID: testID,
		client: client,
	}
	h.testQueue = append(h.testQueue, newItem)
	client.addTestToQueue((len(h.testQueue)-1), newItem.SequenceName, testID)
	h.testQueueUpdateFeed.Send(true)
}

func (h *HttpServer) removeTestFromQueue(testID uuid.UUID) {
	removedTestIndex := 0
	for i := range h.testQueue {
		if (h.testQueue[i].UUID == testID) {
			removedTestIndex = i
			//remove test from client queue
			h.testQueue[i].client.removeTestFromQueue(removedTestIndex)
			//remove test from server queue
			h.testQueue = append(h.testQueue[:i], h.testQueue[i+1:]...)
			break
		}
	}

	if removedTestIndex == len(h.testQueue) {
		return
	}
	// Update client queueNumber after for all tests after removeTestIndex
	for i:=removedTestIndex; i<len(h.testQueue); i++ {
		//update client test queue
		h.testQueue[i].client.updateTestQueue(i)
	}
}

// update client test queues when test finishes
func (h *HttpServer) updateTestQueue() {
	h.removeTestFromQueue(h.testQueue[0].UUID)

	// Check again if the queue still has items before attempting to access them
	if len(h.testQueue) > 0 {
		for i := range h.testQueue {
			h.testQueue[i].client.updateTestQueue(0)
		}
	}
	h.testQueueUpdateFeed.Send(true)
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
			h.updateTestQueue()

		case <-ctx.Done():
			h.l.Info("context done signal received")

			return
		}
	}
}
