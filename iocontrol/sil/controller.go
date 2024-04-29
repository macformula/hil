package sil

import (
	"context"
	"net"

	proto "github.com/macformula/hil/iocontrol/sil/generated"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	_loggerName           = "sil_controller"
	_unknownSignalTypeErr = "unknown signal type"
)

type Controller struct {
	proto.UnimplementedSignalsServer

	l      *zap.Logger
	addr   string
	server *grpc.Server

	digitalSignals map[string]map[string]bool
	analogSignals  map[string]map[string]float64
	signals        map[string]map[string]*proto.SignalInfo
}

// NewController returns a new SIL controller.
func NewController(l *zap.Logger, addr string) *Controller {
	return &Controller{
		l:              l.Named(_loggerName),
		addr:           addr,
		server:         grpc.NewServer(),
		digitalSignals: make(map[string]map[string]bool),
		analogSignals:  make(map[string]map[string]float64),
		signals:        make(map[string]map[string]*proto.SignalInfo),
	}
}

// Open configures the gRPC server to communicate with firmware.
func (c *Controller) Open() error {
	c.l.Info("opening sil controller")
	proto.RegisterSignalsServer(c.server, c)

	listener, err := net.Listen("tcp", c.addr)
	if err != nil {
		return errors.Wrap(err, "listen")
	}

	go func() {
		if err := c.server.Serve(listener); err != nil {
			c.l.Error("grpc serve", zap.Error(err))
		}
	}()

	return nil
}

// Close stops the gRPC server.
func (c *Controller) Close() {
	c.l.Info("closing sil controller")
	c.server.GracefulStop()
}

// registerSignalIfNotExists ensures that any signals that are used by the SIL controller are registered.
// It also prevents any accesses to nil maps.
func (c *Controller) registerSignalIfNotExists(info *proto.SignalInfo) {
	if c.signals[info.EcuName] == nil {
		c.signals[info.EcuName] = make(map[string]*proto.SignalInfo)
	}
	c.signals[info.EcuName][info.SignalName] = info

	switch info.SignalType {
	case proto.SignalType_SIGNAL_TYPE_DIGITAL:
		if c.digitalSignals[info.EcuName] == nil {
			c.digitalSignals[info.EcuName] = make(map[string]bool)
		}
	case proto.SignalType_SIGNAL_TYPE_ADC:
		if c.analogSignals[info.EcuName] == nil {
			c.analogSignals[info.EcuName] = make(map[string]float64)
		}
	}
}

// SetDigital sets an output digital pin for a SIL digital pin.
func (c *Controller) SetDigital(pin *DigitalPin, b bool) {
	c.registerSignalIfNotExists(pin.Info)
	c.digitalSignals[pin.Info.EcuName][pin.Info.SignalName] = b
}

// ReadDigital returns the level of a SIL digital pin.
func (c *Controller) ReadDigital(pin *DigitalPin) bool {
	c.registerSignalIfNotExists(pin.Info)
	return c.digitalSignals[pin.Info.EcuName][pin.Info.SignalName]
}

// WriteVoltage sets the voltage of a SIL analog pin.
func (c *Controller) WriteVoltage(pin *AnalogPin, voltage float64) {
	c.registerSignalIfNotExists(pin.Info)
	c.analogSignals[pin.Info.SignalName][pin.Info.SignalName] = voltage
}

// ReadVoltage returns the voltage of a SIL analog pin.
func (c *Controller) ReadVoltage(pin *AnalogPin) float64 {
	c.registerSignalIfNotExists(pin.Info)
	return c.analogSignals[pin.Info.EcuName][pin.Info.SignalName]
}

// WriteCurrent sets the current of a SIL analog pin (unimplemented for SIL).
func (c *Controller) WriteCurrent(output *AnalogPin, current float64) error {
	return errors.New("unimplemented function on sil controller")
}

// ReadCurrent returns the current of a SIL analog pin (unimplemented for SIL).
func (c *Controller) ReadCurrent(output *AnalogPin) (float64, error) {
	return 0.00, errors.New("unimplemented function on sil controller")
}

// EnumerateRegisteredSignals is a gRPC server call that returns all registered signals.
func (c *Controller) EnumerateRegisteredSignals(_ context.Context, _ *proto.EnumerateRegisteredSignalsRequest) (*proto.EnumerateRegisteredSignalsResponse, error) {
	var registeredSignals []*proto.SignalInfo
	for _, signals := range c.signals {
		for _, info := range signals {
			registeredSignals = append(registeredSignals, info)
		}
	}

	return &proto.EnumerateRegisteredSignalsResponse{
		Status:  true,
		Error:   "",
		Signals: registeredSignals,
	}, nil
}

// WriteSignal is a gRPC server call that sets a signal value depending on the signal type.
func (c *Controller) WriteSignal(_ context.Context, in *proto.WriteSignalRequest) (*proto.WriteSignalResponse, error) {
	info := c.signals[in.EcuName][in.SignalName]

	switch info.SignalType {
	case proto.SignalType_SIGNAL_TYPE_DIGITAL:
		c.SetDigital(&DigitalPin{info}, in.GetValueDigital().GetLevel())
	case proto.SignalType_SIGNAL_TYPE_ADC:
		c.WriteVoltage(&AnalogPin{info}, in.GetValueAdc().GetVoltage())
	default:
		return &proto.WriteSignalResponse{
			Status: false,
			Error:  _unknownSignalTypeErr,
		}, errors.New(_unknownSignalTypeErr)
	}

	return &proto.WriteSignalResponse{
		Status: true,
		Error:  "",
	}, nil
}

// ReadSignal is a gRPC server call that reads a signal value depending on the signal type.
func (c *Controller) ReadSignal(_ context.Context, in *proto.ReadSignalRequest) (*proto.ReadSignalResponse, error) {
	info := c.signals[in.EcuName][in.SignalName]

	switch info.SignalType {
	case proto.SignalType_SIGNAL_TYPE_DIGITAL:
		return &proto.ReadSignalResponse{
			Status: true,
			Error:  "",
			Value: &proto.ReadSignalResponse_ValueDigital{
				ValueDigital: &proto.DigitalSignal{Level: c.ReadDigital(&DigitalPin{info})},
			},
		}, nil
	case proto.SignalType_SIGNAL_TYPE_ADC:
		return &proto.ReadSignalResponse{
			Status: true,
			Error:  "",
			Value: &proto.ReadSignalResponse_ValueAdc{
				ValueAdc: &proto.AdcSignal{Voltage: c.ReadVoltage(&AnalogPin{info})},
			},
		}, nil
	default:
		return &proto.ReadSignalResponse{
			Status: false,
			Error:  _unknownSignalTypeErr,
		}, errors.New(_unknownSignalTypeErr)
	}
}

// RegisterSignal is a gRPC server call that creates an entry in signals for a new signal sent by firmware.
// This allows us to register signals from either the HIL or the firmware.
func (c *Controller) RegisterSignal(_ context.Context, in *proto.RegisterSignalRequest) (*proto.RegisterSignalResponse, error) {
	info := &proto.SignalInfo{
		EcuName:      in.EcuName,
		SignalName:   in.SignalName,
		SignalType:   proto.SignalType_SIGNAL_TYPE_UNKNOWN,
		SignalAccess: in.SignalAccess,
	}

	if in.GetInitialValueDigital() != nil {
		info.SignalType = proto.SignalType_SIGNAL_TYPE_DIGITAL
		c.SetDigital(&DigitalPin{info}, in.GetInitialValueDigital().GetLevel())
	} else if in.GetInitialValueAdc() != nil {
		info.SignalType = proto.SignalType_SIGNAL_TYPE_ADC
		c.WriteVoltage(&AnalogPin{info}, in.GetInitialValueAdc().GetVoltage())
	} else if in.GetInitialValuePwm() != nil {
		info.SignalType = proto.SignalType_SIGNAL_TYPE_PWM
	} else {
		return &proto.RegisterSignalResponse{
			Status: false,
			Error:  _unknownSignalTypeErr,
		}, errors.New(_unknownSignalTypeErr)
	}

	c.registerSignalIfNotExists(info)

	return &proto.RegisterSignalResponse{
		Status: true,
		Error:  "",
	}, nil
}
