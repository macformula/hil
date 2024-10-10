package macformula

import (
	"time"

	"go.uber.org/zap"

	"github.com/macformula/hil/macformula/pinout"
	"github.com/pkg/errors"
)

const (
	_testbenchLoggerName     = "testbench"
	_hvilOnMinVoltage        = 4
	_hvilOffMaxVoltage       = 1
	_defaultEcuOnMinVoltage  = 3
	_defaultEcuOffMinVoltage = 1
)

// TestBench allows control of the test bench in hatch.
type TestBench struct {
	l             *zap.Logger
	pinController *pinout.Controller

	checkLvControllerVoltage bool
	ecuOnMinVoltage          float64
	ecuOffMaxVoltage         float64
	checkHvilFeedback        bool
}

type TestBenchOption func(*TestBench)

func WithCheckLvControllerVoltageOnPowerCycle() TestBenchOption {
	return func(tb *TestBench) {
		tb.checkLvControllerVoltage = true
	}
}

func WithRequiredEcuVoltages(ecuOnMinVoltage, ecuOffMaxVoltage float64) TestBenchOption {
	return func(tb *TestBench) {
		tb.ecuOffMaxVoltage = ecuOffMaxVoltage
		tb.ecuOnMinVoltage = ecuOnMinVoltage
	}
}

func WithCheckHvilFeedbackOnHvilControl() TestBenchOption {
	return func(tb *TestBench) {
		tb.checkHvilFeedback = true
	}
}

func NewTestBench(pc *pinout.Controller, l *zap.Logger, opts ...TestBenchOption) *TestBench {
	tb := &TestBench{
		l:                        l.Named(_testbenchLoggerName),
		pinController:            pc,
		checkLvControllerVoltage: false,
		ecuOnMinVoltage:          _defaultEcuOnMinVoltage,
		ecuOffMaxVoltage:         _defaultEcuOffMinVoltage,
	}

	for _, o := range opts {
		o(tb)
	}

	return tb
}

func (tb *TestBench) PowerCycle() error {
	var voltage float64

	err := tb.pinController.SetDigitalLevel(pinout.GlvmsDisable, true)
	if err != nil {
		return errors.Wrap(err, "set digital level")
	}

	time.Sleep(1 * time.Second)

	if tb.checkLvControllerVoltage {
		voltage, err = tb.pinController.ReadVoltage(pinout.LvController3v3RefVoltage)
		if err != nil {
			return errors.Wrap(err, "read voltage")
		}

		if voltage > tb.ecuOffMaxVoltage {
			return errors.Errorf("lv controller voltage exceeds max voltage for testbench off: max (%v), got (%v)",
				tb.ecuOffMaxVoltage,
				voltage)
		}
	}

	err = tb.pinController.SetDigitalLevel(pinout.GlvmsDisable, false)
	if err != nil {
		return errors.Wrap(err, "set digital level")
	}

	if tb.checkLvControllerVoltage {
		voltage, err = tb.pinController.ReadVoltage(pinout.LvController3v3RefVoltage)
		if err != nil {
			return errors.Wrap(err, "read voltage")
		}

		if voltage < tb.ecuOnMinVoltage {
			return errors.Errorf("lv controller voltage less than min voltage for testbench on: min (%v), got (%v)",
				tb.ecuOnMinVoltage,
				voltage)
		}
	}

	return nil
}

func (tb *TestBench) IsLvControllerEnabled() (bool, error) {
	voltage, err := tb.pinController.ReadVoltage(pinout.LvController3v3RefVoltage)
	if err != nil {
		return false, errors.Wrap(err, "read voltage")
	}

	if voltage < tb.ecuOffMaxVoltage {
		return false, nil
	}

	if voltage > tb.ecuOnMinVoltage {
		return true, nil
	}

	return false, errors.Errorf("lv controller voltage is between the max off value and the min on value (%v<%v<%v)",
		tb.ecuOffMaxVoltage, voltage, tb.ecuOnMinVoltage)
}

func (tb *TestBench) IsFrontControllerEnabled() (bool, error) {
	voltage, err := tb.pinController.ReadVoltage(pinout.FrontController3v3RefVoltage)
	if err != nil {
		return false, errors.Wrap(err, "read voltage")
	}

	if voltage < tb.ecuOffMaxVoltage {
		return false, nil
	}

	if voltage > tb.ecuOnMinVoltage {
		return true, nil
	}

	return false, errors.Errorf("front controller voltage is between the max off value and the min on value (%v<%v<%v)",
		tb.ecuOffMaxVoltage, voltage, tb.ecuOnMinVoltage)
}

func (tb *TestBench) BreakHvil() error {
	var voltage float64

	err := tb.pinController.SetDigitalLevel(pinout.HvilDisable, true)
	if err != nil {
		return errors.Wrap(err, "set digital level")
	}

	if tb.checkHvilFeedback {
		voltage, err = tb.pinController.ReadVoltage(pinout.HvilFeedback)
		if err != nil {
			return errors.Wrap(err, "read voltage")
		}

		if voltage > _hvilOffMaxVoltage {
			return errors.Errorf("hvil voltage exceeds max on hvil break (max: %v, got %v)",
				_hvilOffMaxVoltage, voltage)
		}
	}

	return nil
}

func (tb *TestBench) ResetHvil() error {
	var voltage float64

	err := tb.pinController.SetDigitalLevel(pinout.HvilDisable, false)
	if err != nil {
		return errors.Wrap(err, "set digital level")
	}

	if tb.checkHvilFeedback {
		voltage, err = tb.pinController.ReadVoltage(pinout.HvilFeedback)
		if err != nil {
			return errors.Wrap(err, "read voltage")
		}

		if voltage > _hvilOnMinVoltage {
			return errors.Errorf("hvil voltage less than min on hvil reset (min: %v, got %v)",
				_hvilOnMinVoltage, voltage)
		}
	}

	return nil
}
