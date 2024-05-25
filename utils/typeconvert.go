package utils

// TODO: make this a generic
func BoolToNumeric(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
