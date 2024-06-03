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

	successfulPowerCycle                    bool
	tsalEnabled                             bool
	tsalEnableDuration                      time.Duration
	raspiEnabled                            bool
	raspiEnableDuration                     time.Duration
	frontControllerEnabled                  bool
	frontControllerEnableDuration           time.Duration
	speedgoatEnabled                        bool
	speedgoatEnableDuration                 time.Duration
	accumulatorEnabled                      bool
	accumulatorEnableDuration               time.Duration
	motorPrechargeEnabled                   bool
	motorPrechargeEnableDuration            time.Duration
	motorControllerEnabled                  bool
	motorControllerEnableDuration           time.Duration
	motorPrechageDisabled                   bool
	motorPrechageDisabledDuration           time.Duration
	shutdownEnabledBeforeHvContactorsCan    bool
	shutdownEnabledBeforeHvContactorsOpen   bool
	shutdownCircuitEnabled                  bool
	dcdcEnabledBeforeClosedContactors       bool
	dcdcEnabledAfterClosedContactors        bool
	shutdownCircuitEnableDuration           time.Duration
	inverterSwitchEnabled                   bool
	inverterSwitchEnabledDuration           time.Duration
	inverterEnabledBeforeHvContactorsCan    bool
	inverterEnabledBeforeHvContactorsClosed bool
	inverterEnabledBeforeCommand            bool
}

func newLvStartup(a *macformula.App, l *zap.Logger) *lvStartup {
	return &lvStartup{
		l:  l,
		a:  a,
		tb: a.TestBench,
		lv: a.LvControllerClient,
		fc: a.FrontControllerClient,
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
		return errors.Wrap(err, "poll read tsal on")
	}

	l.raspiEnabled, l.raspiEnableDuration, err = utils.Poll(ctx, l.lv.ReadRaspiOn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read raspi on")
	}

	l.frontControllerEnabled, l.frontControllerEnableDuration, err =
		utils.Poll(ctx, l.lv.ReadFrontControllerOn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read front controller on")
	}

	l.speedgoatEnabled, l.speedgoatEnableDuration, err =
		utils.Poll(ctx, l.lv.ReadSpeedgoatOn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read speedgoat on")
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

	inverterOn, err := l.lv.ReadInverterSwitchOn()
	if err != nil {
		return errors.Wrap(err, "read inverter switch on")
	}

	if inverterOn {
		l.inverterEnabledBeforeHvContactorsCan = true
	}

	shutdownOn, err = l.sendContactorCommandCheckShutdownOn(
		ctx, _hvPositiveOn, _hvNegativeOn, _prechargeOff, 2*time.Second, 100*time.Millisecond)
	if err != nil {
		return errors.Wrap(err, "send contactor command check shutdown on")
	}

	if shutdownOn {
		l.shutdownEnabledBeforeHvContactorsOpen = true
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
		l.inverterEnabledBeforeHvContactorsClosed = true
	}

	dcdcEnabled, err := l.lv.ReadDcdcOn()
	if err != nil {
		return errors.Wrap(err, "read dcdc on")
	}

	if dcdcEnabled {
		l.dcdcEnabledBeforeClosedContactors = true
	}

	err = l.fc.CommandContactors(ctx, _hvPositiveOn, _hvNegativeOn, _prechargeOff)
	if err != nil {
		return errors.Wrap(err, "command contactors")
	}

	time.Sleep(1 * time.Second)

	l.dcdcEnabledAfterClosedContactors, err = l.lv.ReadDcdcOn()
	if err != nil {
		return errors.Wrap(err, "read dcdc on")
	}

	err = l.lv.SetDcdcValid(true)
	if err != nil {
		return errors.Wrap(err, "set dcdc valid")
	}

	inverterOn, err = l.lv.ReadInverterSwitchOn()
	if err != nil {
		return errors.Wrap(err, "read inverter switch on")
	}

	if inverterOn {
		l.inverterEnabledBeforeCommand = true
	}

	err = l.fc.CommandInverter(ctx, _inverterEnable)
	if err != nil {
		return errors.Wrap(err, "command inverter")
	}

	l.inverterSwitchEnabled, l.inverterSwitchEnabledDuration, err = utils.Poll(
		ctx, l.lv.ReadInverterSwitchOn, _lvStartupPollTimeout)
	if err != nil {
		return errors.Wrap(err, "poll read inverter switch on")
	}

	return nil
}

func (l *lvStartup) GetResults() map[flow.Tag]any {
	return map[flow.Tag]any{
		config.LvStartupTags.PowerCycledTestBench:                              l.successfulPowerCycle,
		config.LvStartupTags.TsalEnabled:                                       l.tsalEnabled,
		config.LvStartupTags.TsalTimeToEnableMs:                                int(l.tsalEnableDuration.Milliseconds()),
		config.LvStartupTags.RaspiEnabled:                                      l.raspiEnabled,
		config.LvStartupTags.RaspiTimeToEnableMs:                               int(l.raspiEnableDuration.Milliseconds()),
		config.LvStartupTags.FrontControllerEnabled:                            l.frontControllerEnabled,
		config.LvStartupTags.FrontControllerTimeToEnableMs:                     int(l.frontControllerEnableDuration.Milliseconds()),
		config.LvStartupTags.SpeedgoatEnabled:                                  l.speedgoatEnabled,
		config.LvStartupTags.SpeedgoatTimeToEnableMs:                           int(l.speedgoatEnableDuration.Milliseconds()),
		config.LvStartupTags.AccumulatorEnabled:                                l.accumulatorEnabled,
		config.LvStartupTags.AccumulatorTimeToEnableMs:                         int(l.accumulatorEnableDuration.Milliseconds()),
		config.LvStartupTags.MotorPrechageEnabled:                              l.motorPrechargeEnabled,
		config.LvStartupTags.MotorPrechargeTimeToEnableMs:                      int(l.motorPrechargeEnableDuration.Milliseconds()),
		config.LvStartupTags.MotorControllerEnabled:                            l.motorControllerEnabled,
		config.LvStartupTags.MotorControllerTimeToEnable:                       int(l.motorControllerEnableDuration.Milliseconds()),
		config.LvStartupTags.ShutdownCircuitEnabledBeforeCan:                   l.shutdownEnabledBeforeHvContactorsCan,
		config.LvStartupTags.ShutdownCircuitEnabledBeforeOpenContactors:        l.shutdownEnabledBeforeHvContactorsOpen,
		config.LvStartupTags.ShutdownCircuitEnabled:                            l.shutdownCircuitEnabled,
		config.LvStartupTags.ShutdownCircuitTimeToEnable:                       int(l.shutdownCircuitEnableDuration.Milliseconds()),
		config.LvStartupTags.DcdcEnabledBeforeContactorsClosed:                 l.dcdcEnabledBeforeClosedContactors,
		config.LvStartupTags.DcdcEnabledAfterContactorsClosed:                  l.dcdcEnabledAfterClosedContactors,
		config.LvStartupTags.InverterSwitchEnabledBeforeCan:                    l.inverterEnabledBeforeHvContactorsCan,
		config.LvStartupTags.InverterSwitchEnabledBeforeClosedContactors:       l.inverterEnabledBeforeHvContactorsClosed,
		config.LvStartupTags.InverterSwitchEnabledBeforeFrontControllerCommand: l.inverterEnabledBeforeCommand,
		config.LvStartupTags.InverterSwitchEnabled:                             l.inverterSwitchEnabled,
		config.LvStartupTags.InverterSwitchTimeToEnable:                        int(l.inverterSwitchEnabledDuration.Milliseconds()),
	}
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
