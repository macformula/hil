package writer

type WriterIface interface {
	CreateTraceFile(traceDir string, busName string) error
	WriteFrameToFile(frame *TimestampedFrame) error
	CloseTraceFile() error
}