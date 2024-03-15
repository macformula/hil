package httpdispatcher

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/macformula/hil/orchestrator"
)

type TestQueueItem struct {
	queueNumber int
	UUID        uuid.UUID
}

type Client struct {
	conn      *websocket.Conn
	testQueue []TestQueueItem
	status    chan orchestrator.StatusSignal
	results   chan orchestrator.ResultsSignal
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		// testQueue: make(chan TestQueueItem, 10),
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

// func (c *Client) updateTests() {
// 	for i := range c.testQueue {
// 		if c.testQueue[i].queueNumber > 0 {
// 			c.testQueue[i].queueNumber -= 1
// 		} else {
// 			c.removeTest(i)
// 		}
// 	}
// }
