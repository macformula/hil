package macformula

import (
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/macformula/config"
	"github.com/macformula/hil/macformula/pinout"
)

type App struct {
	Config       *config.Config
	VehCanTracer *canlink.Tracer
	PtCanTracer  *canlink.Tracer

	PinoutController *pinout.Controller
	CurrProcess      *ProcessInfo
}
