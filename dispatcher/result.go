package dispatcher

import "time"

type result struct {
	duration time.Duration
	desc     string
	passed   bool
	name     string
}
