package sequencer

import "context"

type State interface {
	Name() string
	Start(ctx context.Context) error
	FatalError() error
}
