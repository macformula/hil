package iocontrol

// DigitalPin defines methods for digital pin control
type DigitalPin interface {
	Open() error
	SetDirection(Direction) error
	Read() (Level, error)
	Write(Level) error
}
