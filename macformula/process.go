package macformula

// ProcessInfo represents the information about the current process.
type ProcessInfo struct {
	boardSerial string // Serial number of the board (this is a placeholder for now to show how this will be used)
}

// NewProcessInfo creates a new ProcessInfo instance.
func NewProcessInfo() *ProcessInfo {
	return &ProcessInfo{}
}
