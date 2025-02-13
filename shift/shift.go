package main

import "time"

// Each shift from the dataset.json file, validate input JSON against this struct
type Shift struct {
	ShiftID    int64     `json:"ShiftID"`
	EmployeeID int64     `json:"EmployeeID"`
	StartTime  time.Time `json:"StartTime"`
	EndTime    time.Time `json:"EndTime"`
}
