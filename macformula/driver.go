package macformula

import "context"

type Driver interface {
	// StartButton gives access to the start button on the dash.
	StartButton(ctx context.Context, press bool) error
	// Accelerate allows to modify either one or two pedal positions (one for each potentiometer).
	Accelerate(ctx context.Context, pedalPct ...float64) error
	// Brake allows to slow down the vehicle.
	Brake(ctx context.Context, pedalPct float64) error
	// Steer allows setting the current steering wheel angle.
	Steer(ctx context.Context, steeringAngle float64) error
}
