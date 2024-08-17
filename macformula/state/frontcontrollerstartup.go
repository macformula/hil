package state

import (
	"context"
	"github.com/macformula/hil/flow"
	"time"
)

const (
	_frontControllerStartup = "front_controller_startup"
)

type frontControllerStartup struct{}

func (f *frontControllerStartup) Name() string {
	return _frontControllerStartup
}

func (f *frontControllerStartup) Setup(ctx context.Context) error {
	//TODO: indicate start of simulation to front controller

	return nil
}

func (f *frontControllerStartup) Run(ctx context.Context) error {
	panic("implement me")
}

func (f *frontControllerStartup) GetResults() map[flow.Tag]any {
	//TODO implement me
	panic("implement me")
}

func (f *frontControllerStartup) ContinueOnFail() bool {
	//TODO implement me
	panic("implement me")
}

func (f *frontControllerStartup) Timeout() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (f *frontControllerStartup) FatalError() error {
	//TODO implement me
	panic("implement me")
}
