package canlink

// TraceFile is a file that should be able to log CAN frames.
type TraceFile interface {
	// dumpToFile takes a list of CAN frames and writes them to a file.
	dumpToFile(frames []TimestampedFrame, traceDir, busName string) error
}
