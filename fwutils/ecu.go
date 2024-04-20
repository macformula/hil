package fwutils

//go:generate enumer -type=Ecu ecu.go
type Ecu int

const (
	UnknownEcu Ecu = iota
	FrontController
	LvController
	TmsController
)
