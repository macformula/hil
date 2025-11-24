package state

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula"
	"github.com/macformula/hil/macformula/config"
	"github.com/macformula/hil/macformula/ecu/frontcontroller"
	"github.com/macformula/hil/macformula/ecu/lvcontroller"
	"github.com/macformula/hil/macformula/pinout"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
)

const (
	_lvStartupName           = "lv_startup"
	_lvStartupPollTimeout    = 10 * time.Second
	_lvStartupContinueOnFail = true
	_lvStartupTimeout        = 2 * time.Minute

	_hvPositiveOn   = true
	_hvPositiveOff  = false
	_hvNegativeOn   = true
	_hvNegativeOff  = false
	_prechargeOff   = false
	_inverterEnable = true
)

type lvStartup struct {
	l  *zap.Logger
	a  *macformula.App
	tb *macformula.TestBench
	lv *lvcontroller.Client
	fc *frontcontroller.Client

	fatalErr utils.ResettableError

	results map[flow.Tag]any
}

func newLvStartup(a *macformula.App, l *zap.Logger) *lvStartup {
	return &lvStartup{
		l:       l,
		a:       a,
		tb:      a.TestBench,
		lv:      a.LvControllerClient,
		fc:      a.FrontControllerClient,
		results: map[flow.Tag]any{},
	}
}

func (l *lvStartup) Name() string {
	return _lvStartupName
}

func (l *lvStartup) Setup(_ context.Context) error {
	return nil
}

func (l *lvStartup) Run(ctx context.Context) error {
	var (
		r    = l.results
		tags = config.LvStartupTags
	)

	err := l.tb.PowerCycle()
	if err != nil {
		r[tags.PowerCycledTestBench] = false

		return errors.Wrap(err, "power cycle")
	}

	r[tags.PowerCycledTestBench] = true

	r[tags.TsalEnabled], r[tags.TsalTimeToEnableMs], err = pollPinMs(ctx, *l.lv, pinout.TsalEn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read tsal on")
	}

	r[tags.RaspiEnabled], r[tags.RaspiTimeToEnableMs], err = pollPinMs(ctx, *l.lv, pinout.RaspiEn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read raspi on")
	}

	r[tags.FrontControllerEnabled], r[tags.FrontControllerTimeToEnableMs], err =
		pollPinMs(ctx, *l.lv, pinout.FrontControllerEn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read front controller on")
	}

	r[tags.SpeedgoatEnabled], r[tags.SpeedgoatTimeToEnableMs], err =
		pollPinMs(ctx, *l.lv, pinout.SpeedgoatEn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read speedgoat on")
	}

	r[tags.AccumulatorEnabled], r[tags.AccumulatorTimeToEnableMs], err =
		pollPinMs(ctx, *l.lv, pinout.AccumulatorEn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read accumulator on")
	}

	r[tags.MotorPrechageEnabled], r[tags.MotorPrechargeTimeToEnableMs], err =
		pollPinMs(ctx, *l.lv, pinout.MotorControllerPrechargeEn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read motor controller precharge on")
	}

	r[tags.MotorControllerEnabled], r[tags.MotorControllerTimeToEnable], err =
		pollPinMs(ctx, *l.lv, pinout.MotorControllerEn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read motor controller on")
	}

	r[tags.MotorPrechargeDisabled], r[tags.MotorPrechargeTimeToDisableMs], err =
		pollPinMs(ctx, *l.lv, pinout.MotorControllerPrechargeEn, _lvStartupPollTimeout, utils.CheckForFalse())
	if err != nil {
		return errors.Wrap(err, "poll read motor controller precharge on")
	}

	time.Sleep(2 * time.Second)

	shutdownOn, err := l.lv.ReadDigital(pinout.ShutdownCircuitEn)
	if err != nil {
		return errors.Wrap(err, "read shutdown circuit on")
	}

	if shutdownOn {
		r[tags.ShutdownCircuitEnabledBeforeCan] = true
	} else {
		r[tags.ShutdownCircuitEnabledBeforeCan] = false
	}

	inverterOn, err := l.lv.ReadDigital(pinout.InverterSwitchEn)
	if err != nil {
		return errors.Wrap(err, "read inverter switch on")
	}

	if inverterOn {
		r[tags.InverterSwitchEnabledBeforeCan] = true
	} else {
		r[tags.InverterSwitchEnabledBeforeCan] = false
	}

	shutdownOn, err = l.sendContactorCommandCheckShutdownOn(
		ctx, _hvPositiveOn, _hvNegativeOn, _prechargeOff, 2*time.Second, 100*time.Millisecond)
	if err != nil {
		return errors.Wrap(err, "send contactor command check shutdown on")
	}

	if shutdownOn {
		r[tags.ShutdownCircuitEnabledBeforeOpenContactors] = true
	} else {
		r[tags.ShutdownCircuitEnabledBeforeOpenContactors] = false
	}

	r[tags.ShutdownCircuitEnabled], err = l.sendContactorCommandCheckShutdownOn(
		ctx, _hvPositiveOff, _hvNegativeOff, _prechargeOff, 2*time.Second, 100*time.Millisecond)
	if err != nil {
		return errors.Wrap(err, "send contactor command check shutdown on")
	}

	inverterOn, err = l.lv.ReadDigital(pinout.InverterSwitchEn)
	if err != nil {
		return errors.Wrap(err, "read inverter switch on")
	}

	if inverterOn {
		r[tags.InverterSwitchEnabledBeforeClosedContactors] = true
	} else {
		r[tags.InverterSwitchEnabledBeforeClosedContactors] = false
	}

	dcdcEnabled, err := l.lv.ReadDigital(pinout.DcdcEn)
	if err != nil {
		return errors.Wrap(err, "read dcdc on")
	}

	if dcdcEnabled {
		r[tags.DcdcEnabledBeforeContactorsClosed] = true
	} else {
		r[tags.DcdcEnabledBeforeContactorsClosed] = false
	}

	err = l.fc.CommandContactors(ctx, _hvPositiveOn, _hvNegativeOn, _prechargeOff)
	if err != nil {
		return errors.Wrap(err, "command contactors")
	}

	time.Sleep(1 * time.Second)

	r[tags.DcdcEnabledAfterContactorsClosed], err = l.lv.ReadDigital(pinout.DcdcEn)
	if err != nil {
		return errors.Wrap(err, "read dcdc on")
	}

	err = l.lv.SetDcdcValid(true)
	if err != nil {
		return errors.Wrap(err, "set dcdc valid")
	}

	inverterOn, err = l.lv.ReadDigital(pinout.InverterSwitchEn)
	if err != nil {
		return errors.Wrap(err, "read inverter switch on")
	}

	if inverterOn {
		r[tags.InverterSwitchEnabledBeforeFrontControllerCommand] = true
	} else {
		r[tags.InverterSwitchEnabledBeforeFrontControllerCommand] = false
	}

	err = l.fc.CommandInverter(ctx, _inverterEnable)
	if err != nil {
		return errors.Wrap(err, "command inverter")
	}

	r[tags.InverterSwitchEnabled], r[tags.InverterSwitchTimeToEnable], err = pollPinMs(
		ctx, *l.lv, pinout.InverterSwitchEn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read inverter switch on")
	}

	return nil
}

func (l *lvStartup) GetResults() map[flow.Tag]any {
	return l.results
}

func (l *lvStartup) ContinueOnFail() bool {
	return _lvStartupContinueOnFail
}

func (l *lvStartup) Timeout() time.Duration {
	return _lvStartupTimeout
}

func (l *lvStartup) FatalError() error {
	return l.fatalErr.Err()
}

func (l *lvStartup) sendContactorCommandCheckShutdownOn(ctx context.Context,
	hvPositive, hvNegative, precharge bool, timeout time.Duration, period time.Duration) (bool, error) {
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-timeoutChan:
			return false, nil
		case <-ticker.C:
			err := l.fc.CommandContactors(ctx, hvPositive, hvNegative, precharge)
			if err != nil {
				return false, errors.Wrap(err, "command contactors")
			}

			lvl, err := l.lv.ReadDigital(pinout.ShutdownCircuitEn)
			if err != nil {
				return false, errors.Wrap(err, "read shutdown circuit on")
			}

			if lvl {
				return true, nil
			}
		}
	}
}

// pollMs polls the checkFunc until it returns true or the timeout is reached. It wraps the utils.Poll function
// as the results cannot be time.Duration values (so it converts these to int.
func pollMs(ctx context.Context,
	checkFunc utils.CheckFunc,
	timeout time.Duration,
	opts ...utils.PollOption) (bool, int, error) {
	valid, duration, err := utils.Poll(ctx, checkFunc, timeout, opts...)
	return valid, int(duration.Milliseconds()), errors.Wrap(err, "poll")
}

func pollPinMs(ctx context.Context,
	client lvcontroller.Client,
	pin pinout.PhysicalIo,
	timeout time.Duration,
	opts ...utils.PollOption) (bool, int, error) {

	checkFunc := func() (bool, error) {
		return client.ReadDigital(pin)
	}

	return pollMs(ctx, checkFunc, timeout, opts...)
}
