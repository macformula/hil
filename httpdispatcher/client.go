package httpdispatcher

import (
	"github.com/gorilla/websocket"
)

// WebSocket client connection
type Client struct {
	Conn           *websocket.Conn // The WebSocket connection
	Send           chan []byte     // Channel for outgoing messages
	receiveHandler func([]byte)    // Handler for processing received messages
	errchannel     chan error
}

func NewClient(conn *websocket.Conn, receiveHandler func([]byte)) *Client {
	return &Client{
		Conn:           conn,
		Send:           make(chan []byte, 256), // Buffered channel for outgoing messages
		receiveHandler: receiveHandler,
	}
}
