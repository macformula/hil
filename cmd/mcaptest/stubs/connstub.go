package stubs

import (
	"context"
	"errors"
	"net"
	"time"
)

// Define a stub struct that implements net.Conn
type ConnStub struct{}

// Implement the net.Conn interface methods for the stub

func (c *ConnStub) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (c *ConnStub) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (c *ConnStub) Close() error {
	return nil
}

func (c *ConnStub) LocalAddr() net.Addr {
	return &net.IPAddr{IP: net.IPv4(127, 0, 0, 1), Zone: ""}
}

func (c *ConnStub) RemoteAddr() net.Addr {
	return &net.IPAddr{IP: net.IPv4(127, 0, 0, 1), Zone: ""}
}

func (c *ConnStub) SetDeadline(t time.Time) error {
	return nil
}

func (c *ConnStub) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *ConnStub) SetWriteDeadline(t time.Time) error {
	return nil
}

// Define a stub function to replace socketcan.DialContext
func DialContextStub(ctx context.Context, network, address string) (net.Conn, error) {
	// You can customize the behavior of this stub as needed for testing
	// For example, return a ConnStub and a nil error
	if network != "fail" && address != "fail" {
		return &ConnStub{}, nil
	}
	return nil, errors.New("failed to dial context")
}
