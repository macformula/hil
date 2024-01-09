package iocontrol

//go:generate enumer -type=Level types.go
type Level int

const (
	Unknown Level = iota
	Low
	High
)

//go:generate enumer -type=Direction types.go
type Direction int

const (
	Input Direction = iota
	Output
)
