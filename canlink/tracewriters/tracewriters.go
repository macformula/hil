package tracewriters

type TraceWriter interface {
	CreateTraceFile(traceDir string, busName string) error
	WriteFrameToFile(frame *TimestampedFrame) error
	CloseTraceFile() error
}