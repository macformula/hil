package httpdispatcher

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/macformula/hil/orchestrator"
)

type ClientTestQueueItem struct {
	SequenceName	string
	QueueIndex 		int
	UUID        	uuid.UUID
}

type Client struct {
	conn      			*websocket.Conn
	testQueue 		 	[]ClientTestQueueItem
	testQueueUpdated 	chan bool
	status    		 	chan orchestrator.StatusSignal
	results   		 	chan orchestrator.ResultsSignal
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: 				conn,
		testQueue: 			make([]ClientTestQueueItem, 0),
		testQueueUpdated: 	make(chan bool),
		status:    			make(chan orchestrator.StatusSignal),
		results:   			make(chan orchestrator.ResultsSignal),
	}
}

func (c *Client) addTestToQueue(queueIndex int, sequenceName string,testID uuid.UUID) {
	//add queuenumber and testID to testQueue
	newItem := ClientTestQueueItem{
		SequenceName: 	sequenceName,
		QueueIndex: 	queueIndex,
		UUID:       	testID,
	}
	c.testQueue = append(c.testQueue, newItem)
}

func (c *Client) removeTestFromQueue(queueIndexRemoved int) {
	for i:=0; i<len(c.testQueue); i++ {
		if c.testQueue[i].QueueIndex == queueIndexRemoved {
			c.testQueue = append(c.testQueue[:i], c.testQueue[i+1:]...)
		}
	}
}

func (c *Client) updateTestQueue(queueIndexRemoved int) {
	if len(c.testQueue) == 0 {
		return
	}
	// lower test Index by 1
	for i := range c.testQueue {
		if c.testQueue[i].QueueIndex == queueIndexRemoved + 1 {
			c.testQueue[i].QueueIndex -= 1
			break
		}
	}
}