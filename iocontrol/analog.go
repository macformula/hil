package iocontrol

// AnalogPin is a single analog input/output
type AnalogPin interface {
	String() string
	IsAnalogPin()
}
