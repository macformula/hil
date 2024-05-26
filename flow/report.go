package flow

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID           uuid.UUID `json:"id"`
	SequenceName string    `json:"sequence_name"`
	DateTime     time.Time `json:"date_time"`
	Description  string    `json:"description"`
}

type Test struct {
	ID            uuid.UUID `json:"id"`
	Description   string    `json:"description"`
	Value         bool      `json:"value"`
	CompOp        string    `json:"comp_op"`
	UpperLimit    string    `json:"upper_limit"`
	LowerLimit    string    `json:"lower_limit"`
	ExpectedValue string    `json:"expected_val"`
	Type          bool      `json:"type"`
	Unit          string    `json:"unit"`
}
