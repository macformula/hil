package iocontrol

// DigitalPin is a single digital input/output
type DigitalPin interface {
	String() string
	IsDigitalPin()
}
