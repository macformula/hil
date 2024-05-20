package utils

type Numeric interface {
	(int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64)
}

// TODO: make this a generic

func BoolToNumeric(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
