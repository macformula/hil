package hil

type ProcessInfo struct {
	boardSerial string
}

func NewProcessInfo() *ProcessInfo {
	return &ProcessInfo{}
}
