package dispatcher

import (
	"context"
)

type Dispatcher interface {
	Start(context.Context) error
	GetStartSignal(context.Context) <-chan struct{}
	GetResults(context.Context) error
	Stop(context.Context) error
}
