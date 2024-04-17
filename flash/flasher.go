package flash

type FlasherIface interface {
	String() string
	Flash(string) error
	Open() error
	Close() error
}
