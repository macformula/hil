package httpdispatcher

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type TestQueueItem struct {
	queueNumber int
	UUID        uuid.UUID
}

type Client struct {
	conn      *websocket.Conn
	testQueue []TestQueueItem
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		// testQueue: make(chan TestQueueItem, 10),
		testQueue: make([]TestQueueItem, 0),
	}
}

func (c *Client) addTestToQueue(conn *websocket.Conn, queueNumber int, testID uuid.UUID) {
	c.conn = conn
	//add queuenumber and testID to testQueue
	newItem := TestQueueItem{
		queueNumber: queueNumber,
		UUID:        testID,
	}
	c.testQueue = append(c.testQueue, newItem)
}
