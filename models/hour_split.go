package models

import (
	"time"
)

// Hours store the start and end time of a shift and the hours worked
type Hours struct {
	Start time.Time
	End   time.Time
	Hours float64
}
