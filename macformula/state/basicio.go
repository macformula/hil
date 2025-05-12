package state

import (
	"context"
	"fmt"
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
	_BasicIoName           = "test_state"
	_BasicIoPollTimeout    = 10 * time.Second
	_BasicIoContinueOnFail = true
	_BasicIoTimeout        = 2 * time.Minute
)

type BasicIo struct {
	l *zap.Logger
	a *macformula.App

	fatalErr utils.ResettableError

	results map[flow.Tag]any
}

// This state was designed to test the BasicIO project in firmware directory
func newBasicIo(a *macformula.App, l *zap.Logger) *BasicIo {
	return &BasicIo{
		l:       l,
		a:       a,
		results: map[flow.Tag]any{},
	}
}

func (l *BasicIo) Name() string {
	return _BasicIoName
}

func (l *BasicIo) Setup(_ context.Context) error {
	return nil
}

func (l *BasicIo) Run(ctx context.Context) error {
	l.l.Info("Starting BasicIO sil test")

	l.a.PinModel.RegisterDigitalInput("DemoProject", "IndicatorLed")
	l.a.PinModel.RegisterDigitalOutput("DemoProject", "IndicatorButton")

	var (
		r    = l.results
		tags = config.BasicIoTags
	)

	l.a.PinoutController.SetDigitalLevel(pinout.IndicatorButton, true)

	time.Sleep(200 * time.Millisecond)

	ledLevel, err := l.a.PinoutController.ReadDigitalLevel(pinout.IndicatorLed)
	if err != nil {
		return errors.Wrap(err, "read digital level (indicator led)")
	}

	l.l.Info(fmt.Sprintf("Indicator led is %t", ledLevel))

	r[tags.LedMatchesButtonHigh] = ledLevel

	l.a.PinoutController.SetDigitalLevel(pinout.IndicatorButton, false)

	time.Sleep(100 * time.Millisecond)

	ledLevel, err = l.a.PinoutController.ReadDigitalLevel(pinout.IndicatorLed)
	if err != nil {
		return errors.Wrap(err, "read digital level (indicator led)")
	}
	r[tags.LedMatchesButtonLow] = !ledLevel

	return nil
}

func (l *BasicIo) GetResults() map[flow.Tag]any {
	return l.results
}

func (l *BasicIo) ContinueOnFail() bool {
	return _BasicIoContinueOnFail
}

func (l *BasicIo) Timeout() time.Duration {
	return _BasicIoTimeout
}

func (l *BasicIo) FatalError() error {
	return l.fatalErr.Err()
}
