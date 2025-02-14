package models

import "time"

// JSON Report of each employees hours worked in a week.
type EmployeeSummary struct {
	EmployeeID int64 `json:"EmployeeID"`
	// This should be a date formated as "YYYY-MM-DD"
	StartOfWeek   time.Time
	RegularHours  float64 `json:"RegularHours"`
	OverTimeHours float64 `json:"OvertimeHours"`
	// Array of ShiftIDs that overlap in the same week
	InvalidShifts []int64 `json:"InvalidShifts"`
}
