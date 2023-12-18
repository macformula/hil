package orchestrator

import (
	"io"
)

type Dispatcher interface {
	io.Closer
	Start() <-chan struct{}
	Handshake(<-chan struct{}) <-chan struct{}
	Quit() <-chan struct{}
}
