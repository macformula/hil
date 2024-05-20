package state

import (
	"context"
	"github.com/macformula/hil/macformula/ecu/frontcontroller"
	"time"

	"go.uber.org/zap"

	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula"
	"github.com/macformula/hil/macformula/ecu/lvcontroller"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
)

const (
	_lvStartupName           = "lv_startup"
	_lvStartupPollTimeout    = 10 * time.Second
	_lvStartupContinueOnFail = true
	_lvStartupTimeout        = 1 * time.Minute

	_hvPositiveOn  = true
	_hvPositiveOff = false
	_hvNegativeOn  = true
	_hvNegativeOff = true
	_prechargeOff  = true
)

type lvStartup struct {
	l  *zap.Logger
	a  *macformula.App
	tb *macformula.TestBench
	lv *lvcontroller.Client
	fc *frontcontroller.Client

	fatalErr utils.ResettableError

	successfulPowerCycle                  bool
	tsalEnabled                           bool
	tsalEnableDuration                    time.Duration
	accumulatorEnabled                    bool
	accumulatorEnableDuration             time.Duration
	motorPrechargeEnabled                 bool
	motorPrechargeEnableDuration          time.Duration
	motorControllerEnabled                bool
	motorControllerEnableDuration         time.Duration
	motorPrechageDisabled                 bool
	motorPrechageDisabledDuration         time.Duration
	shutdownEnabledBeforeHvContactorsCan  bool
	shutdownEnabledBeforeHvContactorsOpen bool
	shutdownCircuitEnabled                bool
	shutdownCircuitEnsableDuration        time.Duration
	inverterSwitchEnabled                 bool
	inverterSwitchEnabledDuration         time.Duration
	inverterEnabledBeforeHvContactorsCan  bool
	inverterEnabledBeforeHvContactorsOpen bool
	inverterEnabledBeforeCommand          bool
}

func newLvStartup(a *macformula.App, l *zap.Logger) *lvStartup {
	return &lvStartup{
		l:  l,
		a:  a,
		tb: a.TestBench,
		lv: a.LvControllerClient,
	}
}

func (l *lvStartup) Name() string {
	return _lvStartupName
}

func (l *lvStartup) Setup(_ context.Context) error {
	return nil
}

func (l *lvStartup) Run(ctx context.Context) error {
	err := l.tb.PowerCycle()
	if err != nil {
		l.successfulPowerCycle = false

		return errors.Wrap(err, "power cycle")
	}

	l.successfulPowerCycle = true

	l.tsalEnabled, l.tsalEnableDuration, err = utils.Poll(ctx, l.lv.ReadTsalEn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read tsal en")
	}

	l.accumulatorEnabled, l.accumulatorEnableDuration, err =
		utils.Poll(ctx, l.lv.ReadAccumulatorOn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read accumulator on")
	}

	l.motorPrechargeEnabled, l.motorPrechargeEnableDuration, err =
		utils.Poll(ctx, l.lv.ReadMotorControllerPrechargeOn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read motor controller precharge on")
	}

	l.motorControllerEnabled, l.motorControllerEnableDuration, err =
		utils.Poll(ctx, l.lv.ReadMotorControllerOn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read motor controller on")
	}

	l.motorPrechageDisabled, l.motorPrechageDisabledDuration, err =
		utils.Poll(ctx, l.lv.ReadMotorControllerPrechargeOff, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read motor controller precharge on")
	}

	// TODO: make const
	time.Sleep(2 * time.Second)

	shutdownOn, err := l.lv.ReadShutdownCircuitOn()
	if err != nil {
		return errors.Wrap(err, "read shutdown circuit on")
	}

	if shutdownOn {
		l.shutdownEnabledBeforeHvContactorsCan = true
	}

	shutdownOn, err = l.sendContactorCommandCheckShutdownOn(
		ctx, _hvPositiveOn, _hvNegativeOn, _prechargeOff, 2*time.Second, 100*time.Millisecond)
	if err != nil {
		return errors.Wrap(err, "send contactor command check shutdown on")
	}

	if shutdownOn {
		l.shutdownEnabledBeforeHvContactorsOpen = true
	}

	inverterOn, err := l.lv.ReadInverterSwitchOn()
	if err != nil {
		return errors.Wrap(err, "read inverter switch on")
	}

	if inverterOn {
		l.inverterEnabledBeforeHvContactorsCan = true
	}

	l.shutdownCircuitEnabled, err = l.sendContactorCommandCheckShutdownOn(
		ctx, _hvPositiveOff, _hvNegativeOff, _prechargeOff, 2*time.Second, 100*time.Millisecond)
	if err != nil {
		return errors.Wrap(err, "send contactor command check shutdown on")
	}

	inverterOn, err = l.lv.ReadInverterSwitchOn()
	if err != nil {
		return errors.Wrap(err, "read inverter switch on")
	}

	if inverterOn {
		l.inverterEnabledBeforeHvContactorsOpen = true
	}

	err = l.fc.CommandInverter(ctx, true)
	if err != nil {
		return errors.Wrap(err, "command inverter on")
	}

	l.inverterSwitchEnabled, l.inverterSwitchEnabledDuration, err = utils.Poll(
		ctx, l.lv.ReadInverterSwitchOn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read inverter switch on")
	}

}

func (l *lvStartup) GetResults() map[flow.Tag]any {
	//TODO implement me
	panic("implement me")
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
	hvpositive, hvNegative, precharge bool, timeout time.Duration, period time.Duration) (bool, error) {
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
			err := l.fc.CommandContactors(ctx, hvpositive, hvNegative, precharge)
			if err != nil {
				return false, errors.Wrap(err, "command contactors")
			}

			lvl, err := l.lv.ReadShutdownCircuitOn()
			if err != nil {
				return false, errors.Wrap(err, "read shutdown circuit on")
			}

			if lvl {
				return true, nil
			}
		}
	}
}
