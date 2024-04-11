package httpdispatcher

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
)

type ClientTestQueueItem struct {
	QueueIndex 	int
	Sequence	flow.Sequence
	UUID        uuid.UUID
}

type Client struct {
	conn      *websocket.Conn
	testQueue []ClientTestQueueItem
	status    chan orchestrator.StatusSignal
	results   chan orchestrator.ResultsSignal
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		testQueue: make([]ClientTestQueueItem, 0),
		status:    make(chan orchestrator.StatusSignal),
		results:   make(chan orchestrator.ResultsSignal),
	}
}


// FIX THIS SHIT
func (c *Client) addTestToQueue(queueIndex int, sequence flow.Sequence,testID uuid.UUID) {
	//add queuenumber and testID to testQueue
	newItem := ClientTestQueueItem{
		QueueIndex: queueIndex,
		Sequence: 	sequence,
		UUID:       testID,
	}
	c.testQueue = append(c.testQueue, newItem)
}

func (c *Client) removeTestFromQueue(queueIndexRemoved int) {
	for i:=0; i<len(c.testQueue); i++ {
		if c.testQueue[i].queueIndex == queueIndexRemoved {
			c.testQueue = append(c.testQueue[:i], c.testQueue[i+1:]...)
		}
	}
}

func (c *Client) updateTestFromQueue(queueIndexRemoved int) {
	// lower test Index by 1
	for i := range c.testQueue {
		if c.testQueue[i].queueIndex == queueIndexRemoved + 1 {
			c.testQueue[i].queueIndex -= 1
			break
		}
	}
}

// FIX THIS SHIT
// func (c *Client) updateTests() {
// 	select {
// 	case <-c.results:
// 		// Update queue position of client tests
// 		for i := range c.testQueue {
// 			if c.testQueue[i].queueNumber > 0 {
// 				c.testQueue[i].queueNumber -= 1
// 			} else {
// 				c.removeTest(i)
// 			}
// 		}
// 		queueData, _ := json.Marshal(c.testQueue)
// 		c.conn.WriteMessage(websocket.TextMessage, queueData)
// 	}
// }