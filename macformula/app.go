package macformula

import (
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/macformula/config"
	"github.com/macformula/hil/macformula/ecu/frontcontroller"
	"github.com/macformula/hil/macformula/ecu/lvcontroller"
	"github.com/macformula/hil/macformula/pinout"
)

type App struct {
	Config                *config.Config
	VehCanTracer          *canlink.Tracer
	PtCanTracer           *canlink.Tracer
	PinoutController      *pinout.Controller
	TestBench             *TestBench
	LvControllerClient    *lvcontroller.Client
	FrontControllerClient *frontcontroller.Client
	VehCanClient          *canlink.CanClient
	PtCanClient           *canlink.CanClient

	CurrProcess *ProcessInfo
}
