package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Requirements:
// ShiftID is 64 bit integer and required
// EmployeeID is 64 bit integer and required
// StartTime is a time.Time object localized in the CST timezone,and RFC3339
// EndTime is a time.Time object localized in the CST timezone,and RFC3339
// RawStartTime is a string used to validate format of StartTime
// RawEndTime is a string used to validate format of EndTime
type Shift struct {
	ShiftID      int64 `json:"ShiftID" validate:"required"`
	EmployeeID   int64 `json:"EmployeeID" validate:"required"`
	StartTime    time.Time
	EndTime      time.Time
	RawStartTime string `json:"StartTime" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	RawEndTime   string `json:"EndTime" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// Use go-playground/validatorv/v10 to validate the Shift struct fields
func (s *Shift) ValidateShift() error {
	validate := validator.New()
	return validate.Struct(s)
}
