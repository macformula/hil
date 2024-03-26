package main

import (
	"context"
	"fmt"
	"time"

	"github.com/macformula/hil/can_gen/VEH_CAN"
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/config"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
)

const (
	_configPath = "/opt/macfe/etc/config.yaml"
	_busName    = "VEH_CAN"
	_fastMsgs   = 15
	_slowMsgs   = 5
)

func main() {
	conf, err := config.NewConfig(_configPath)
	if err != nil {
		panic("error reading conf file")
	}

	ctx := context.Background()

	cfg := zap.NewDevelopmentConfig()
	formattedTime := time.Now().Format("2006.01.02_15.04.05")
	fileName := fmt.Sprintf("/opt/macfe/traces/logs/tracetest_%s.log", formattedTime)
	cfg.OutputPaths = []string{fileName}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	conn, err := socketcan.DialContext(context.Background(), "can", "can0")
	if err != nil {
		panic(err)
	}

	canClient := canlink.NewCANClient(VEH_CAN.Messages(), conn)

	m, err := canlink.NewMcap(logger, canClient, conf.TracerDirectory, conf.BusName)
	if err != nil {
		panic(err)
	}
	a, err := canlink.NewAsc(logger, conf.TracerDirectory, conf.BusName)
	if err != nil {
		panic(err)
	}

	tracer := canlink.NewTracer(
		conf.CanInterface,
		conf.TracerDirectory,
		logger,
		conn,
		canlink.WithBusName(_busName),
		canlink.WithFiles(m, a))

	canClient.Open()
	err = tracer.Open(ctx)
	if err != nil {
		logger.Error("open tracer", zap.Error(err))
		return
	}

	err = tracer.StartTrace(ctx)
	if err != nil {
		logger.Error("start trace", zap.Error(err))
	}

	println("Start First Test")
	// First Test

	for i := 0; i <= _fastMsgs; i++ {
		p := VEH_CAN.NewPack_State()
		p.SetAvg_Cell_Voltage(float64(i))
		p.SetPopulated_Cells(uint8(i * 2))
		p.SetPack_Inst_Voltage(float64(i * 3))
		p.SetPack_Current(float64(i * 4))

		newChannel2 := VEH_CAN.NewPack_SOC()
		newChannel2.SetMaximum_Pack_Voltage(float64(i))
		newChannel2.SetPack_SOC(float64(i))

		println("Sent faster: ", i)
		frame, err := p.MarshalFrame()
		if err != nil {
			return
		}
		frame2, err := newChannel2.MarshalFrame()
		if err != nil {
			return
		}

		time.Sleep(500 * time.Millisecond)
		canClient.Send(context.Background(), frame)
		canClient.Send(context.Background(), frame2)
	}

	time.Sleep(1 * time.Second)

	for i := 0; i < _slowMsgs; i++ {
		println("Sent slower: ", i)
		c := can.Frame{
			ID:         1600,
			Length:     8,
			Data:       can.Data{byte(i)},
			IsRemote:   false,
			IsExtended: false,
		}
		time.Sleep(1000 * time.Millisecond)
		canClient.Send(context.Background(), c)
	}

	time.Sleep(500 * time.Millisecond)

	println("End First Test")

	err = tracer.StopTrace()
	if err != nil {
		logger.Error("stop trace", zap.Error(err))
	}

	err = tracer.Close()
	if err != nil {
		logger.Error("close tracer", zap.Error(err))
	}

	err = canClient.Close()
	if err != nil {
		logger.Error("close canclient", zap.Error(err))
	}

	if tracer.Error() != nil {
		logger.Error("tracer error", zap.Error(tracer.Error()))
	}

	logger.Info("End of Main")
}
