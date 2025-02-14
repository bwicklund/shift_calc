package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Each shift from the dataset.json file, validate input JSON against this struct
type Shift struct {
	ShiftID      int64 `json:"ShiftID" validate:"required"`
	EmployeeID   int64 `json:"EmployeeID" validate:"required"`
	StartTime    time.Time
	EndTime      time.Time
	RawStartTime string `json:"StartTime" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	RawEndTime   string `json:"EndTime" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

func (s *Shift) ValidateShift() error {
	validate := validator.New()
	return validate.Struct(s)
}
