package commander

import (
	"context"
	"github.com/macformula/hilmatic/sequencer"
	"go.uber.org/zap"
)

type Commander struct {
	l *zap.Logger
	s *sequencer.Sequencer
}

func NewCommander(l *zap.Logger) *Commander {
	return &Commander{
		l: l,
	}
}

func (c *Commander) Open(ctx context.Context) error {

}

func (c *Commander) StartTest(ctx context.Context) error {

}

func (c *Commander) Close(ctx context.Context) error {

}
