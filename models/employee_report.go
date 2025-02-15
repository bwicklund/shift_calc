package models

// JSON Report of each employees hours worked in a week.
// Requirements:
// EmployeeID is 64 bit integer
// StartOfWeek is a string in the format "YYYY-MM-DD"
// RegularHours is a float64
// OverTimeHours is a float64
// InvalidShifts is an array of 64 bit integers (ShiftIDs)

// If I have time I should validate this output... TODO
type EmployeeReport struct {
	EmployeeID    int64   `json:"EmployeeID"`
	StartOfWeek   string  `json:"StartOfWeek"`
	RegularHours  float64 `json:"RegularHours"`
	OverTimeHours float64 `json:"OvertimeHours"`
	InvalidShifts []int64 `json:"InvalidShifts"`
}
