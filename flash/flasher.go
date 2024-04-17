package flash

type FlasherIface interface {
	String() string
	Flash(string) error
	Connect() error
	Disconnect() error
}
