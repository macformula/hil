package sil

import (
	"context"
	"fmt"
	"io"

	"net"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	flatbuffers "github.com/google/flatbuffers/go"
	pb "github.com/macformula/hil/iocontrol/sil/proto"
	signals "github.com/macformula/hil/iocontrol/sil/signals"
)

const (
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/signals.proto
type FbController struct {
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

// NewFbController returns a new SIL FbController.
func NewFbController(port int, l *zap.Logger) *FbController {
	return &FbController{
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
func (c *FbController) Open(ctx context.Context) error {
	c.l.Info("opening sil FbController")

	addr := fmt.Sprintf("localhost:%v", c.port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		c.l.Error(fmt.Sprintf("creating listener: %s", errors.Wrap(err, "creating listener")))
	}
	defer listener.Close();

	c.l.Info(fmt.Sprintf("sil listening on %s", addr))

	for {
		conn, err := listener.Accept()
		if err != nil {
			c.l.Fatal(fmt.Sprintf("accepting sil client: %s", err))
		}
		go c.handleConnection(conn)
	}
}

// Close stops the gRPC server.
func (c *FbController) Close() error {
	c.l.Info("closing sil FbController")
	c.server.GracefulStop()

	return nil
}

func (c *FbController) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 2024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
            if err != io.EOF {
                c.l.Error(fmt.Sprintf("read error: %s", err))
            }
            break
        }
		request := signals.GetRootAsRequest(buffer, 0)
		
		unionTable := new(flatbuffers.Table)
		if request.Request(unionTable) {
			requestType := request.RequestType()
			if requestType == signals.RequestTypeReadRequest {
				unionRequest := new(signals.ReadRequest)
				unionRequest.Init(unionTable.Bytes, unionTable.Pos)

				ecu := string(unionRequest.EcuName())
				sig_name := string(unionRequest.SignalName())
				c.l.Info(fmt.Sprintf("the ecu is: %s, sig is %s", ecu, sig_name))
			}

		}
		

		// switch requestType {
		// case signals.RequestTypeReadRequest:
		// 	c.l.Info("read request")
		// case signals.RequestTypeSetRequest:
		// 	c.l.Info("set request")
		// case signals.RequestTypeRegisterRequest:
		// 	c.l.Info("register request")
		// default:
		// 	c.l.Info("none matched")
		// }

		fmt.Printf("recevied %x", buffer[:n])
	}
}

func (c *FbController) registerDigitalOutput(ecuName, sigName string) {
	mapSet(c.digitalOutputPins, ecuName, sigName, false)
	mapSet(c.signalInfo, ecuName, sigName, &pb.SignalInfo{
		EcuName:         ecuName,
		SignalName:      sigName,
		SignalType:      pb.SignalType_SIGNAL_TYPE_DIGITAL,
		SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_INPUT, // Purposely reversed
	})
}

func (c *FbController) registerDigitalInput(ecuName, sigName string) {
	mapSet(c.digitalInputPins, ecuName, sigName, false)
	mapSet(c.signalInfo, ecuName, sigName, &pb.SignalInfo{
		EcuName:         ecuName,
		SignalName:      sigName,
		SignalType:      pb.SignalType_SIGNAL_TYPE_DIGITAL,
		SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_OUTPUT, // Purposely reversed
	})
}

func (c *FbController) registerAnalogOutput(ecuName, sigName string) {
	mapSet(c.analogOutputPins, ecuName, sigName, 0.0)
	mapSet(c.signalInfo, ecuName, sigName, &pb.SignalInfo{
		EcuName:         ecuName,
		SignalName:      sigName,
		SignalType:      pb.SignalType_SIGNAL_TYPE_ANALOG,
		SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_INPUT, // Purposely reversed
	})
}

func (c *FbController) registerAnalogInput(ecuName, sigName string) {
	mapSet(c.analogInputPins, ecuName, sigName, 0.0)
	mapSet(c.signalInfo, ecuName, sigName, &pb.SignalInfo{
		EcuName:         ecuName,
		SignalName:      sigName,
		SignalType:      pb.SignalType_SIGNAL_TYPE_ANALOG,
		SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_OUTPUT, // Purposely reversed
	})
}

// SetDigital sets an output digital pin for a SIL digital pin.
func (c *FbController) SetDigital(pin *DigitalPin, level bool) error {
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
func (c *FbController) ReadDigital(pin *DigitalPin) (bool, error) {
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
func (c *FbController) WriteVoltage(pin *AnalogPin, voltage float64) error {
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
func (c *FbController) ReadVoltage(pin *AnalogPin) (float64, error) {
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
func (c *FbController) WriteCurrent(_ *AnalogPin, _ float64) error {
	return errors.New("unimplemented function on sil FbController")
}

// ReadCurrent returns the current of a SIL analog pin (unimplemented for SIL).
func (c *FbController) ReadCurrent(_ *AnalogPin) (float64, error) {
	return 0.00, errors.New("unimplemented function on sil FbController")
}

// EnumerateRegisteredSignals is a gRPC server call that returns all registered signalInfo.
func (c *FbController) EnumerateRegisteredSignals(
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
func (c *FbController) WriteSignal(_ context.Context, in *pb.WriteSignalRequest) (*pb.WriteSignalResponse, error) {
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
func (c *FbController) ReadSignal(_ context.Context, in *pb.ReadSignalRequest) (*pb.ReadSignalResponse, error) {
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

func (c *FbController) RegisterSignal(_ context.Context, in *pb.RegisterSignalRequest) (*pb.RegisterSignalResponse, error) {
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

