package canlink

type Handler interface {
	Name() string
	Handle(chan TimestampedFrame) error
}
