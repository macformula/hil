package iocontrol

type DigitalPin interface {
	SetDirection(Direction) error
	Read() (Level, error)
	Write(Level) error
}
