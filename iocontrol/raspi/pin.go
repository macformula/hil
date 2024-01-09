package raspi

type Pin uint8

//go:generate enumer -type=Pin pin.go
const (
	Gpio0 Pin = iota
	Gpio1
	Gpio2
	Gpio3
	Gpio4
	Gpio5
	Gpio6
	Gpio7
	Gpio8
	Gpio9
	Gpio10
	Gpio11
	Gpio12
	Gpio13
	Gpio14
	Gpio15
	Gpio16
	Gpio17
	Gpio18
	Gpio19
	Gpio20
	Gpio21
	Gpio22
	Gpio23
	Gpio24
	Gpio25
	Gpio26
	Gpio27
)
