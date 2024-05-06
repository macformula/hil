package fwutils

// FlasherIface defines methods for connecting to a target and flashing it with firmware
type FlasherIface interface {
	// String gets the type of flasher
	String() string
	// Flash loads a binary on the target
	Flash(binName string) error
	// Connect establishes a connection with a target
	Connect(ecu Ecu) error
	// Disconnect closes the connection with the target
	Disconnect() error
}
