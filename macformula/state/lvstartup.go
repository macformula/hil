package state

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula"
)

const (
	_lvStartupName = "lv_startup"
)

type lvStartup struct {
	l *zap.Logger
	a *macformula.App
}

func newLvStartup(a *macformula.App, l *zap.Logger) *lvStartup {
	return &lvStartup{
		l: l,
		a: a,
	}
}

func (l *lvStartup) Name() string {
	return _lvStartupName
}

func (l *lvStartup) Setup(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (l *lvStartup) Run(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (l *lvStartup) GetResults() map[flow.Tag]any {
	//TODO implement me
	panic("implement me")
}

func (l *lvStartup) ContinueOnFail() bool {
	//TODO implement me
	panic("implement me")
}

func (l *lvStartup) Timeout() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (l *lvStartup) FatalError() error {
	//TODO implement me
	panic("implement me")
}
