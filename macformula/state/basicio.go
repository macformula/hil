package state

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula"
	"github.com/macformula/hil/macformula/config"
	"github.com/macformula/hil/macformula/pinout"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
)

const (
	_demoStateName           = "demo_state"
	_demoStatePollTimeout    = 10 * time.Second
	_demoStateContinueOnFail = true
	_demoStateTimeout        = 2 * time.Minute
	_sleepTime = 1 * time.Second
)

type demoState struct {
	l  *zap.Logger
	a  *macformula.App
	p  *pinout.Controller

	fatalErr utils.ResettableError
	results map[flow.Tag]any
}

func newDemoState(a *macformula.App, l *zap.Logger) *demoState {
	return &demoState{
		l:       l,
		a:       a,
		p: 			a.PinoutController,
		results: map[flow.Tag]any{},
	}
}

func (l *demoState) Name() string {
	return _demoStateName
}

func (l *demoState) Setup(_ context.Context) error {
	return nil
}

func (l *demoState) Run(ctx context.Context) error {
	var (
		r    = l.results
		tags = config.DemoStateTags
	)

	err := l.p.SetDigitalLevel(pinout.IndicatorButton, true)
	if err != nil {
		return errors.Wrap(err, "set digital level (indicator button)")
	}
	time.Sleep(_sleepTime)
	indicatorLedOn, err := l.p.ReadDigitalLevel(pinout.IndicatorLed)
	if err != nil {
		return errors.Wrap(err, "read digital level (indicator led)")
	}

	if indicatorLedOn {
		r[tags.MirrorsButtonDown] = true
	} else {
		r[tags.MirrorsButtonDown] = false
	}

	err = l.p.SetDigitalLevel(pinout.IndicatorButton, false)
	if err != nil {
		return errors.Wrap(err, "set digital level (indicator button)")
	}
	time.Sleep(_sleepTime)
	indicatorLedOn, err = l.p.ReadDigitalLevel(pinout.IndicatorLed)
	if err != nil {
		return errors.Wrap(err, "read digital level (indicator led)")
	}

	if !indicatorLedOn {
		r[tags.MirrorsButtonUp] = true
	} else {
		r[tags.MirrorsButtonUp] = false
	}

	return nil
}

func (l *demoState) GetResults() map[flow.Tag]any {
	return l.results
}

func (l *demoState) ContinueOnFail() bool {
	return _demoStateContinueOnFail
}

func (l *demoState) Timeout() time.Duration {
	return _demoStateTimeout
}

func (l *demoState) FatalError() error {
	return l.fatalErr.Err()
}
