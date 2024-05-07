package utils

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type CheckFunc func() (bool, error)

func Poll(ctx context.Context, checkFunc CheckFunc, timeout time.Duration) (bool, time.Duration, error) {
	start := time.Now()
	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ctx.Done():
			return false, time.Since(start), ctx.Err()
		case <-timeoutChan:
			return false, time.Since(start), nil
		default:
			checkValid, err := checkFunc()
			if err != nil {
				return false, time.Since(start), errors.Wrap(err, "check func")
			}

			if checkValid {
				return true, time.Since(start), nil
			}
		}
	}
}

func Sleep(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}
