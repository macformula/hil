package sil

import (
	"context"
	"fmt"
	"io"

	"net"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	flatbuffers "github.com/google/flatbuffers/go"
	pb "github.com/macformula/hil/iocontrol/sil/proto"
	signals "github.com/macformula/hil/iocontrol/sil/signals"
)

const ()

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/signals.proto
type FbController struct {
	pb.UnimplementedSignalsServer

	l      *zap.Logger
	port   int
	server *grpc.Server
	pins   *PinModel
}

// NewFbController returns a new SIL FbController.
func NewFbController(port int, l *zap.Logger) *FbController {
	return &FbController{
		l:    l.Named(_loggerName),
		port: port,
		pins: NewPinModel(),
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
	defer listener.Close()

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
		_, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				c.l.Error(fmt.Sprintf("read error: %s", err))
			}
			break
		}
		request := signals.GetRootAsRequest(buffer, 0)

		unionTable := new(flatbuffers.Table)
		if request.Request(unionTable) {
			c.l.Info("recvieved request")
			requestType := request.RequestType()
			switch requestType {
			case signals.RequestTypeReadRequest:
				unionRequest := new(signals.ReadRequest)
				unionRequest.Init(unionTable.Bytes, unionTable.Pos)

				ecu := string(unionRequest.EcuName())
				sig_name := string(unionRequest.SignalName())
				sig_type := unionRequest.SignalType()
				sig_direction := unionRequest.SignalDirection()

				switch sig_type {
				case signals.SIGNAL_TYPEDIGITAL:
					switch sig_direction {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.ReadDigitalInput(ecu, sig_name)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.ReadDigitalOutput(ecu, sig_name)
					}
				case signals.SIGNAL_TYPEANALOG:
					switch sig_direction {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.ReadAnalogInput(ecu, sig_name)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.ReadAnalogOutput(ecu, sig_name)
					}
				}
			case signals.RequestTypeSetRequest:
				unionRequest := new(signals.SetRequest)
				unionRequest.Init(unionTable.Bytes, unionTable.Pos)

				ecu := string(unionRequest.EcuName())
				sig_name := string(unionRequest.SignalName())
				sig_type := unionRequest.SignalType()
				sig_direction := unionRequest.SignalDirection()

				unionTable = new(flatbuffers.Table)
				var value bool
				var voltage float64
				if unionRequest.SignalValue(unionTable) {
					switch unionRequest.SignalValueType() {
					case signals.SignalValueDigital:
						unionSignalValue := new(signals.Digital)
						unionSignalValue.Init(unionTable.Bytes, unionTable.Pos)

						value = unionSignalValue.Value()

					case signals.SignalValueAnalog:
						unionSignalValue := new(signals.Analog)
						unionSignalValue.Init(unionTable.Bytes, unionTable.Pos)

						voltage = unionSignalValue.Voltage()
					}
				}
				switch sig_type {
				case signals.SIGNAL_TYPEDIGITAL:
					switch sig_direction {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.SetDigitalInput(ecu, sig_name, value)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.SetDigitalOutput(ecu, sig_name, value)
					}
				case signals.SIGNAL_TYPEANALOG:
					switch sig_direction {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.SetAnalogInput(ecu, sig_name, voltage)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.SetAnalogOutput(ecu, sig_name, voltage)
					}
				}

			case signals.RequestTypeRegisterRequest:
				unionRequest := new(signals.RegisterRequest)
				unionRequest.Init(unionTable.Bytes, unionTable.Pos)

				ecu := string(unionRequest.EcuName())
				sig_name := string(unionRequest.SignalName())
				sig_type := unionRequest.SignalType()
				sig_direction := unionRequest.SignalDirection()

				switch sig_type {
				case signals.SIGNAL_TYPEDIGITAL:
					switch sig_direction {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.RegisterDigitalInput(ecu, sig_name)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.RegisterDigitalOutput(ecu, sig_name)
					}
				case signals.SIGNAL_TYPEANALOG:
					switch sig_direction {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.RegisterAnalogInput(ecu, sig_name)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.RegisterAnalogOutput(ecu, sig_name)
					}
				}
			}

		}
	}
}

// WriteCurrent sets the current of a SIL analog pin (unimplemented for SIL).
func (c *FbController) WriteCurrent(_ *AnalogPin, _ float64) error {
	return errors.New("unimplemented function on sil FbController")
}

// ReadCurrent returns the current of a SIL analog pin (unimplemented for SIL).
func (c *FbController) ReadCurrent(_ *AnalogPin) (float64, error) {
	return 0.00, errors.New("unimplemented function on sil FbController")
}
