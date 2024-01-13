package raspi

type Pin uint8

//go:generate enumer -type=Pin pin.go

/* Define GPIOs that are available on your specific machine by
 * adding/removing Pin gpio types.
 * In the case below, pins 7-11, 23, and 25 are reserved for the
 * 2-Channel CAN Hat.
 */
const (
	Gpio0 Pin = iota
	Gpio1
	Gpio2
	Gpio3
	Gpio4
	Gpio5
	Gpio6
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
	Gpio24
	Gpio26
	Gpio27
)
