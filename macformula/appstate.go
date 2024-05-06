package macformula

import (
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/config"
)

type AppState struct {
	Config       *config.Config
	VehCanTracer *canlink.Tracer
	PtCanTracer  *canlink.Tracer

	CurrProcess *ProcessInfo
}
