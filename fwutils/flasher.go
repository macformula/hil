package fwutils

type FlasherIface interface {
	String() string
	Flash(binName string) error
	Connect(ecu Ecu) error
	Disconnect() error
}
