package macformula

import (
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/macformula/config"
	"github.com/macformula/hil/macformula/ecu/frontcontroller"
	"github.com/macformula/hil/macformula/ecu/lvcontroller"
	"github.com/macformula/hil/macformula/pinout"
)

// App represents the main application, it persists accross multiple sequence runs.
type App struct {
	Config                *config.Config
	VehBusManager *canlink.BusManager
	PtBusManager *canlink.BusManager
	VehCanTracer          *canlink.Tracer
	PtCanTracer           *canlink.Tracer
	PinoutController      *pinout.Controller
	TestBench             *TestBench
	LvControllerClient    *lvcontroller.Client
	FrontControllerClient *frontcontroller.Client

	CurrProcess *ProcessInfo
}
