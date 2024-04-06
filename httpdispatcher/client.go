package httpdispatcher

import (
	"encoding/json"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/macformula/hil/orchestrator"
)

type TestQueueItem struct {
	queueNumber int
	UUID        uuid.UUID
}

type Client struct {
	l         *zap.Logger
	conn      *websocket.Conn
	testQueue []TestQueueItem
	status    chan orchestrator.StatusSignal
	results   chan orchestrator.ResultsSignal
}

func NewClient(conn *websocket.Conn, l *zap.Logger) *Client {
	return &Client{
		l:         l.Named(_clientLoggerName),
		conn:      conn,
		testQueue: make([]TestQueueItem, 0),
		status:    make(chan orchestrator.StatusSignal),
		results:   make(chan orchestrator.ResultsSignal),
	}
}

func (c *Client) addTest(queueNumber int, testID uuid.UUID) {
	//add queuenumber and testID to testQueue
	newItem := TestQueueItem{
		queueNumber: queueNumber,
		UUID:        testID,
	}
	c.testQueue = append(c.testQueue, newItem)
}

func (c *Client) removeTest(testIndex int) {
	c.testQueue = append(c.testQueue[:testIndex], c.testQueue[testIndex+1:]...)
}

func (c *Client) updateTests() {
	select {
	case <-c.results:
		c.l.Info("updateTests, result came in")
		// Update queue position of client tests
		for i := range c.testQueue {
			if c.testQueue[i].queueNumber > 0 {
				c.testQueue[i].queueNumber -= 1
			} else {
				c.removeTest(i)
			}
		}
		queueData, _ := json.Marshal(c.testQueue)
		c.conn.WriteMessage(websocket.TextMessage, queueData)
	}
}
