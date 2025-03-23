package sil

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/macformula/hil/iocontrol/sil/proto"
)

const (
	_loggerName           = "sil_controller"
	_unknownSignalTypeErr = "unknown signal type"
	_unknownSignalDirErr  = "unknown signal direction"
	_pwmUnsupportedErr    = "pwm pins are currently unsupported"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/signals.proto
type Controller struct {
	pb.UnimplementedSignalsServer

	l      *zap.Logger
	port   int
	server *grpc.Server

	digitalInputPins  map[string]map[string]bool
	digitalOutputPins map[string]map[string]bool
	analogInputPins   map[string]map[string]float64
	analogOutputPins  map[string]map[string]float64
	signalInfo        map[string]map[string]*pb.SignalInfo

	digitalInputMtx  *sync.Mutex
	digitalOutputMtx *sync.Mutex
	analogInputMtx   *sync.Mutex
	analogOutputMtx  *sync.Mutex
}

// NewController returns a new SIL controller.
func NewController(port int, l *zap.Logger) *Controller {
	return &Controller{
		l:                 l.Named(_loggerName),
		port:              port,
		digitalInputPins:  make(map[string]map[string]bool),
		digitalOutputPins: make(map[string]map[string]bool),
		analogInputPins:   make(map[string]map[string]float64),
		analogOutputPins:  make(map[string]map[string]float64),
		signalInfo:        make(map[string]map[string]*pb.SignalInfo),
		digitalInputMtx:   &sync.Mutex{},
		digitalOutputMtx:  &sync.Mutex{},
		analogInputMtx:    &sync.Mutex{},
		analogOutputMtx:   &sync.Mutex{},
	}
}

// Open configures the gRPC server to communicate with firmware.
func (c *Controller) Open(ctx context.Context) error {
	c.l.Info("opening sil controller")

	addr := fmt.Sprintf("localhost:%v", c.port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "listen")
	}

	c.server = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	pb.RegisterSignalsServer(c.server, c)

	go func() {
		c.l.Info("starting server", zap.String("listening_at", addr))
		err = c.server.Serve(listener)
		if err != nil {
			c.l.Error("grpc serve", zap.Error(err))
		}
	}()

	return nil
}

// Close stops the gRPC server.
func (c *Controller) Close() error {
	c.l.Info("closing sil controller")
	c.server.GracefulStop()

	return nil
}

func (c *Controller) registerDigitalOutput(ecuName, sigName string) {
	mapSet(c.digitalOutputPins, ecuName, sigName, false)
	mapSet(c.signalInfo, ecuName, sigName, &pb.SignalInfo{
		EcuName:         ecuName,
		SignalName:      sigName,
		SignalType:      pb.SignalType_SIGNAL_TYPE_DIGITAL,
		SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_INPUT, // Purposely reversed
	})
}

func (c *Controller) registerDigitalInput(ecuName, sigName string) {
	mapSet(c.digitalInputPins, ecuName, sigName, false)
	mapSet(c.signalInfo, ecuName, sigName, &pb.SignalInfo{
		EcuName:         ecuName,
		SignalName:      sigName,
		SignalType:      pb.SignalType_SIGNAL_TYPE_DIGITAL,
		SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_OUTPUT, // Purposely reversed
	})
}

func (c *Controller) registerAnalogOutput(ecuName, sigName string) {
	mapSet(c.analogOutputPins, ecuName, sigName, 0.0)
	mapSet(c.signalInfo, ecuName, sigName, &pb.SignalInfo{
		EcuName:         ecuName,
		SignalName:      sigName,
		SignalType:      pb.SignalType_SIGNAL_TYPE_ANALOG,
		SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_INPUT, // Purposely reversed
	})
}

func (c *Controller) registerAnalogInput(ecuName, sigName string) {
	mapSet(c.analogInputPins, ecuName, sigName, 0.0)
	mapSet(c.signalInfo, ecuName, sigName, &pb.SignalInfo{
		EcuName:         ecuName,
		SignalName:      sigName,
		SignalType:      pb.SignalType_SIGNAL_TYPE_ANALOG,
		SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_OUTPUT, // Purposely reversed
	})
}

// SetDigital sets an output digital pin for a SIL digital pin.
func (c *Controller) SetDigital(pin *DigitalPin, level bool) error {
	c.digitalOutputMtx.Lock()
	defer c.digitalOutputMtx.Unlock()

	_, ok := mapLookup(c.digitalOutputPins, pin.info.EcuName, pin.info.SignalName)
	if !ok {
		c.registerDigitalOutput(pin.info.EcuName, pin.info.SignalName)
	}

	c.digitalOutputPins[pin.info.EcuName][pin.info.SignalName] = level

	return nil
}

// ReadDigital returns the level of a SIL digital pin.
func (c *Controller) ReadDigital(pin *DigitalPin) (bool, error) {
	c.digitalInputMtx.Lock()
	defer c.digitalInputMtx.Unlock()

	level, ok := mapLookup(c.digitalInputPins, pin.info.EcuName, pin.info.SignalName)
	if !ok {
		return false, errors.Errorf("no entry for ecu name (%s) signal name (%s)",
			pin.info.EcuName, pin.info.SignalName)
	}

	return level, nil
}

// WriteVoltage sets the voltage of a SIL analog pin.
func (c *Controller) WriteVoltage(pin *AnalogPin, voltage float64) error {
	c.analogOutputMtx.Lock()
	defer c.analogOutputMtx.Unlock()

	_, ok := mapLookup(c.analogOutputPins, pin.info.EcuName, pin.info.SignalName)
	if !ok {
		c.registerAnalogOutput(pin.info.EcuName, pin.info.SignalName)
	}

	c.analogOutputPins[pin.info.EcuName][pin.info.SignalName] = voltage

	return nil
}

// ReadVoltage returns the voltage of a SIL analog pin.
func (c *Controller) ReadVoltage(pin *AnalogPin) (float64, error) {
	c.analogInputMtx.Lock()
	defer c.analogInputMtx.Unlock()

	voltage, ok := mapLookup(c.analogInputPins, pin.info.EcuName, pin.info.SignalName)
	if !ok {
		return 0.0, errors.Errorf("no entry for ecu name (%s) signal name (%s)",
			pin.info.EcuName, pin.info.SignalName)
	}

	return voltage, nil
}

// WriteCurrent sets the current of a SIL analog pin (unimplemented for SIL).
func (c *Controller) WriteCurrent(_ *AnalogPin, _ float64) error {
	return errors.New("unimplemented function on sil controller")
}

// ReadCurrent returns the current of a SIL analog pin (unimplemented for SIL).
func (c *Controller) ReadCurrent(_ *AnalogPin) (float64, error) {
	return 0.00, errors.New("unimplemented function on sil controller")
}

// EnumerateRegisteredSignals is a gRPC server call that returns all registered signalInfo.
func (c *Controller) EnumerateRegisteredSignals(
	_ context.Context, _ *pb.EnumerateRegisteredSignalsRequest) (*pb.EnumerateRegisteredSignalsResponse, error) {
	var signals []*pb.SignalInfo

	for _, ecuSignals := range c.signalInfo {
		for _, info := range ecuSignals {
			signals = append(signals, info)
		}
	}

	return &pb.EnumerateRegisteredSignalsResponse{
		Status:  true,
		Error:   "",
		Signals: signals,
	}, nil
}

// WriteSignal is a gRPC server call that sets a signal value depending on the signal type.
func (c *Controller) WriteSignal(_ context.Context, in *pb.WriteSignalRequest) (*pb.WriteSignalResponse, error) {
	switch in.Value.(type) {
	// NOTE: if we get a client request to write a signal, it should go to an input map. The opposite is also true.
	case *pb.WriteSignalRequest_ValueAnalog:
		_, ok := mapLookup(c.analogInputPins, in.EcuName, in.SignalName)
		if !ok {
			return &pb.WriteSignalResponse{
				Status: false,
				Error: fmt.Sprintf("no writeable analog signal with ecu name (%s) and signal name (%s)",
					in.EcuName, in.SignalName),
			}, nil
		}

		c.analogInputMtx.Lock()
		defer c.analogInputMtx.Unlock()

		c.analogInputPins[in.EcuName][in.SignalName] = in.GetValueAnalog().GetVoltage()
	case *pb.WriteSignalRequest_ValueDigital:
		_, ok := mapLookup(c.digitalInputPins, in.EcuName, in.SignalName)
		if !ok {
			return &pb.WriteSignalResponse{
				Status: false,
				Error: fmt.Sprintf("no writeable digital signal with ecu name (%s) and signal name (%s)",
					in.EcuName, in.SignalName),
			}, nil
		}

		c.digitalInputMtx.Lock()
		defer c.digitalInputMtx.Unlock()

		c.digitalInputPins[in.EcuName][in.SignalName] = in.GetValueDigital().GetLevel()
	case *pb.WriteSignalRequest_ValuePwm:
		return &pb.WriteSignalResponse{
			Status: false,
			Error:  _pwmUnsupportedErr,
		}, nil
	default:
		return &pb.WriteSignalResponse{
			Status: false,
			Error:  _unknownSignalTypeErr,
		}, nil
	}

	return &pb.WriteSignalResponse{
		Status: true,
		Error:  "",
	}, nil
}

// ReadSignal is a gRPC server call that reads a signal value depending on the signal type.
func (c *Controller) ReadSignal(_ context.Context, in *pb.ReadSignalRequest) (*pb.ReadSignalResponse, error) {
	switch in.SignalType {
	case pb.SignalType_SIGNAL_TYPE_DIGITAL:
		// Allow reading outputs or inputs. An output signal from the client perspective is an input on the server
		// perspective and vise-versa.

		var pins map[string]map[string]bool
		switch in.SignalDirection {
		case pb.SignalDirection_SIGNAL_DIRECTION_INPUT:
			c.digitalOutputMtx.Lock()
			defer c.digitalOutputMtx.Unlock()
			pins = c.digitalOutputPins
		case pb.SignalDirection_SIGNAL_DIRECTION_OUTPUT:
			c.digitalInputMtx.Lock()
			defer c.digitalInputMtx.Unlock()
			pins = c.digitalInputPins
		default:
			return &pb.ReadSignalResponse{
				Status: false,
				Error:  _unknownSignalDirErr,
				Value:  &pb.ReadSignalResponse_ValueDigital{ValueDigital: &pb.DigitalSignal{Level: false}},
			}, nil

		}

		level, ok := mapLookup(pins, in.EcuName, in.SignalName)
		if !ok {
			return &pb.ReadSignalResponse{
				Status: false,
				Error: fmt.Sprintf("no digital signal with ecu name (%s) and signal name (%s)",
					in.EcuName, in.SignalName),
				Value: &pb.ReadSignalResponse_ValueDigital{ValueDigital: &pb.DigitalSignal{Level: false}},
			}, nil
		}

		return &pb.ReadSignalResponse{
			Status: true,
			Error:  "",
			Value:  &pb.ReadSignalResponse_ValueDigital{ValueDigital: &pb.DigitalSignal{Level: level}},
		}, nil
	case pb.SignalType_SIGNAL_TYPE_ANALOG:
		// Allow reading outputs or inputs. An output signal from the client perspective is an input on the server
		// perspective and vise-versa.

		var pins map[string]map[string]float64
		switch in.SignalDirection {
		case pb.SignalDirection_SIGNAL_DIRECTION_INPUT:
			c.analogOutputMtx.Lock()
			defer c.analogOutputMtx.Unlock()
			pins = c.analogOutputPins
		case pb.SignalDirection_SIGNAL_DIRECTION_OUTPUT:
			c.analogInputMtx.Lock()
			defer c.analogInputMtx.Lock()
			pins = c.analogInputPins
		default:
			return &pb.ReadSignalResponse{
				Status: false,
				Error:  _unknownSignalDirErr,
				Value:  &pb.ReadSignalResponse_ValueAnalog{ValueAnalog: &pb.AnalogSignal{Voltage: 0.0}},
			}, nil
		}

		voltage, ok := mapLookup(pins, in.EcuName, in.SignalName)
		if !ok {
			return &pb.ReadSignalResponse{
				Status: false,
				Error: fmt.Sprintf("no analog signal with ecu name (%s) and signal name (%s)",
					in.EcuName, in.SignalName),
				Value: &pb.ReadSignalResponse_ValueAnalog{ValueAnalog: &pb.AnalogSignal{Voltage: 0.0}},
			}, nil
		}

		return &pb.ReadSignalResponse{
			Status: true,
			Error:  "",
			Value:  &pb.ReadSignalResponse_ValueAnalog{ValueAnalog: &pb.AnalogSignal{Voltage: voltage}},
		}, nil
	default:
		return &pb.ReadSignalResponse{
			Status: false,
			Error:  _unknownSignalTypeErr,
			Value:  &pb.ReadSignalResponse_ValueDigital{ValueDigital: &pb.DigitalSignal{Level: false}},
		}, nil
	}
}

func (c *Controller) RegisterSignal(_ context.Context, in *pb.RegisterSignalRequest) (*pb.RegisterSignalResponse, error) {
	switch in.SignalType {
	case pb.SignalType_SIGNAL_TYPE_DIGITAL:
		// An output signal from the client perspective is an input on the server perspective and vise-versa.
		switch in.SignalDirection {
		case pb.SignalDirection_SIGNAL_DIRECTION_INPUT:
			c.registerDigitalOutput(in.EcuName, in.SignalName)
		case pb.SignalDirection_SIGNAL_DIRECTION_OUTPUT:
			c.registerDigitalInput(in.EcuName, in.SignalName)
		default:
			return &pb.RegisterSignalResponse{
				Status: false,
				Error:  _unknownSignalDirErr,
			}, nil
		}
	case pb.SignalType_SIGNAL_TYPE_ANALOG:
		// An output signal from the client perspective is an input on the server perspective and vise-versa.
		switch in.SignalDirection {
		case pb.SignalDirection_SIGNAL_DIRECTION_INPUT:
			c.registerAnalogOutput(in.EcuName, in.SignalName)
		case pb.SignalDirection_SIGNAL_DIRECTION_OUTPUT:
			c.registerAnalogInput(in.EcuName, in.SignalName)
		default:
			return &pb.RegisterSignalResponse{
				Status: false,
				Error:  _unknownSignalDirErr,
			}, nil
		}
	case pb.SignalType_SIGNAL_TYPE_PWM:
		return &pb.RegisterSignalResponse{
			Status: false,
			Error:  _pwmUnsupportedErr,
		}, nil
	default:
		return &pb.RegisterSignalResponse{
			Status: false,
			Error:  _unknownSignalTypeErr,
		}, nil
	}

	return &pb.RegisterSignalResponse{
		Status: true,
		Error:  "",
	}, nil
}
