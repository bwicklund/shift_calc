package services

import (
	"os"
	"testing"
)

func TestLoadShifts(t *testing.T) {
	// Create a temporary JSON file with test data
	tempFile, err := os.CreateTemp("", "shifts_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Sample shift data
	jsonData := `[
		{"ShiftID": 1, "EmployeeID": 101, "StartTime": "2024-02-14T08:00:00Z", "EndTime": "2024-02-14T16:00:00Z"},
		{"ShiftID": 2, "EmployeeID": 102, "StartTime": "2024-02-14T17:00:00Z", "EndTime": "2024-02-14T23:00:00Z"}
	]`

	_, _ = tempFile.Write([]byte(jsonData))

	// Load shifts from the temp file
	shifts, err := LoadShifts(tempFile.Name())
	if err != nil {
		t.Errorf("LoadShifts returned an error: %v", err)
	}

	// Check if data wasparsed correctly
	if len(shifts) != 2 {
		t.Errorf("Expected 2 shifts, got %d", len(shifts))
	}
	if shifts[0].EmployeeID != 101 {
		t.Errorf("Expected employeeid 101, got %d", shifts[0].EmployeeID)
	}
	if shifts[0].ShiftID != 1 {
		t.Errorf("Expected shiftid 1, got %d", shifts[0].ShiftID)
	}
	// Validate the time is in CST timezone
	if shifts[0].StartTime.String() != "2024-02-14 02:00:00 -0600 CST" {
		t.Errorf("Expected start time 2024-02-14 02:00:00 -0600 CST, got %s", shifts[0].StartTime.String())
	}
	if shifts[0].EndTime.String() != "2024-02-14 10:00:00 -0600 CST" {
		t.Errorf("Expected end time 2024-02-14 10:00:00 -0600 CST, got %s", shifts[0].EndTime.String())
	}
	if shifts[1].EmployeeID != 102 {
		t.Errorf("Expected employeeid 102, got %d", shifts[1].EmployeeID)
	}
	if shifts[1].ShiftID != 2 {
		t.Errorf("Expected shiftid 2, got %d", shifts[1].ShiftID)
	}
	if shifts[1].StartTime.String() != "2024-02-14 11:00:00 -0600 CST" {
		t.Errorf("Expected start time 2024-02-14 11:00:00 -0600 CST, got %s", shifts[1].StartTime.String())
	}
	if shifts[1].EndTime.String() != "2024-02-14 17:00:00 -0600 CST" {
		t.Errorf("Expected end time 2024-02-14 17:00:00 -0600 CST, got %s", shifts[1].EndTime.String())
	}
}
func TestCondenseShiftData(t *testing.T) {
	// Create a temporary JSON file with test data
	tempFile, err := os.CreateTemp("", "shifts_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Sample shift data
	jsonData := `[
			{"ShiftID": 1, "EmployeeID": 101, "StartTime": "2024-02-12T09:00:00Z", "EndTime": "2024-02-12T17:00:00Z"},
			{"ShiftID": 2, "EmployeeID": 101, "StartTime": "2024-02-13T09:00:00Z", "EndTime": "2024-02-13T17:00:00Z"},
			{"ShiftID": 3, "EmployeeID": 101, "StartTime": "2024-02-16T07:00:00Z", "EndTime": "2024-02-16T20:00:00Z"},
			{"ShiftID": 4, "EmployeeID": 101, "StartTime": "2024-02-15T00:00:00Z", "EndTime": "2024-02-15T20:00:00Z"},
			{"ShiftID": 5, "EmployeeID": 101, "StartTime": "2024-02-14T09:00:00Z", "EndTime": "2024-02-14T20:00:00Z"},
			{"ShiftID": 6, "EmployeeID": 101, "StartTime": "2024-02-14T15:00:00Z", "EndTime": "2024-02-14T22:00:00Z"},
			{"ShiftID": 7, "EmployeeID": 102, "StartTime": "2024-02-17T22:00:00Z", "EndTime": "2024-02-18T08:00:00Z"},
			{"ShiftID": 8, "EmployeeID": 101, "StartTime": "2024-03-10T07:30:00Z", "EndTime": "2024-03-10T09:30:00Z"}
		]`

	_, _ = tempFile.Write([]byte(jsonData))

	// Load shifts from the temp file
	shifts, err := LoadShifts(tempFile.Name())
	if err != nil {
		t.Errorf("LoadShifts returned an error: %v", err)
	}

	employeeReports := CondenseShiftData(shifts)

	if len(employeeReports) != 4 {
		t.Errorf("Expected 4 employee reports, got %d", len(employeeReports))
	}

	for _, report := range employeeReports {
		// Regular hours
		if report.EmployeeID == 101 && report.StartOfWeek == "2024-02-11" && report.RegularHours != 40 {
			t.Errorf("Expected 40 RegularHours for EmployeeID 101 week starting 2024-02-11, got %f", report.RegularHours)
		}
		// Overtime hours
		if report.EmployeeID == 101 && report.StartOfWeek == "2024-02-11" && report.OverTimeHours != 9 {
			t.Errorf("Expected 9 OverTimeHours for EmployeeID 101 week starting 2024-02-11, got %f", report.OverTimeHours)
		}
		// Invalid shifts
		if report.EmployeeID == 101 && report.StartOfWeek == "2024-02-11" && len(report.InvalidShifts) != 2 {
			t.Errorf("Expected 2 InvalidShifts for EmployeeID 101 week starting 2024-02-11, got %d", len(report.InvalidShifts))
		}
		// Spanning week boundary
		if report.EmployeeID == 102 && report.StartOfWeek == "2024-02-11" && report.RegularHours != 8 {
			t.Errorf("Expected 8 RegularHours for EmployeeID 102 week starting 2024-02-11, got %f", report.RegularHours)
		}
		if report.EmployeeID == 102 && report.StartOfWeek == "2024-02-18" && report.RegularHours != 2 {
			t.Errorf("Expected 2 RegularHours for EmployeeID 102 week starting 2024-02-18, got %f", report.RegularHours)
		}
		// DST change shift, shift starts at 1:30 AM CST and ends at 4:30 AM CDT ....
		// At 2:00 AM CST, the clock skips forward to 3:00 AM CDT so the shift is only 2 hours
		if report.EmployeeID == 101 && report.StartOfWeek == "2024-03-09" && report.RegularHours != 2 {
			t.Errorf("Expected 2 RegularHours for EmployeeID 101 week starting 2024-03-09, got %f", report.RegularHours)
		}
	}
}
