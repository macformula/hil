package orchestrator

import (
	"context"
	"github.com/macformula/hil/flow"
	"io"
)

type Dispatcher interface {
	io.Closer
	Open(ctx context.Context) error
	Start() <-chan flow.Sequence
	Quit() <-chan struct{}
	RecoverFromFatal() <-chan struct{}
	Progress() chan flow.Progress
}
