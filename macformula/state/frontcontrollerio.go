package state

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula"
	"github.com/macformula/hil/macformula/pinout"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

const _frontControllerIo = "front_controller_io"

type frontControllerIo struct {
	l             *zap.Logger
	a             *macformula.App
	pinController *pinout.Controller
}

func newFrontControllerIo(a *macformula.App, l *zap.Logger) *frontControllerIo {
	return &frontControllerIo{
		l:             l.Named(_frontControllerIo),
		a:             a,
		pinController: a.PinoutController,
	}
}

func (f *frontControllerIo) Name() string {
	return _frontControllerIo
}

func (f *frontControllerIo) Setup(_ context.Context) error {
	err := f.pinController.SetDigitalLevel(pinout.WaitForStart, true)
	if err != nil {
		return errors.Wrap(err, "set digital level")
	}

	time.Sleep(100 * time.Millisecond)

	return nil
}

func (f *frontControllerIo) Run(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (f *frontControllerIo) GetResults() map[flow.Tag]any {
	//TODO implement me
	panic("implement me")
}

func (f *frontControllerIo) ContinueOnFail() bool {
	//TODO implement me
	panic("implement me")
}

func (f *frontControllerIo) Timeout() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (f *frontControllerIo) FatalError() error {
	//TODO implement me
	panic("implement me")
}
