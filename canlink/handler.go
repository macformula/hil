package canlink

// Handler interface describes an acceptable receiver to the bus manager.
type Handler interface {
	Name() string
	Handle(chan TimestampedFrame, chan struct{}) error
}
