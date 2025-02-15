package utils_test

import (
	"os"
	"shift_calc/services"
	"shift_calc/utils"
	"testing"
	"time"
)

func TestLocalizeTime(t *testing.T) {
	// Localize time to CST/CDT timezone, before and after DST change
	tests := []struct {
		input    string
		expected string
	}{
		{"2024-03-10T01:30:00Z", "2024-03-09 19:30:00 -0600 CST"},
		{"2024-03-10T08:30:00Z", "2024-03-10 03:30:00 -0500 CDT"},
	}

	for _, tt := range tests {
		localizedTime, _ := utils.LocalizeTime(tt.input)
		if localizedTime.String() != tt.expected {
			t.Errorf("Expected %q, got %q", tt.expected, localizedTime.String())
		}
	}
}

func TestStartOfWeek(t *testing.T) {
	tests := []struct {
		input    time.Time
		expected string
	}{
		{time.Date(2024, 2, 14, 12, 0, 0, 0, time.UTC), "2024-02-11"},
		{time.Date(2024, 2, 12, 12, 0, 0, 0, time.UTC), "2024-02-11"},
		{time.Date(2024, 2, 18, 12, 0, 0, 0, time.UTC), "2024-02-18"},
	}

	// CST/CDT timezone
	loc, _ := time.LoadLocation("America/Chicago")

	for _, tt := range tests {
		result := utils.StartOfWeek(tt.input.In(loc))
		if result != tt.expected {
			t.Errorf("Expected %q, got %q", tt.expected, result)
		}
	}
}

func TestReturnHoursForShift(t *testing.T) {
	// Create a temporary JSON file with test data
	tempFile, err := os.CreateTemp("", "shifts_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Sample shift data
	jsonData := `[
		{"ShiftID": 1, "EmployeeID": 101, "StartTime": "2024-02-12T09:00:00Z", "EndTime": "2024-02-12T17:00:00Z"},
		{"ShiftID": 2, "EmployeeID": 101, "StartTime": "2024-02-17T22:00:00Z", "EndTime": "2024-02-18T08:00:00Z"}
	]`

	_, _ = tempFile.Write([]byte(jsonData))

	// Load shifts from the temp file
	shifts, err := services.LoadShifts(tempFile.Name())
	if err != nil {
		t.Errorf("LoadShifts returned an error: %v", err)
	}

	currentWeekHours, nextWeekHours := utils.ReturnHoursForShift(shifts[0].StartTime, shifts[0].EndTime)
	if currentWeekHours.Hours != 8 {
		t.Errorf("Expected 8 hours, got %f", currentWeekHours.Hours)
	}
	if nextWeekHours != nil {
		t.Errorf("Expected nil, got %v", nextWeekHours)
	}

	currentWeekHours, nextWeekHours = utils.ReturnHoursForShift(shifts[1].StartTime, shifts[1].EndTime)
	if currentWeekHours.Hours != 8 {
		t.Errorf("Expected 10 hours, got %f", currentWeekHours.Hours)
	}
	if nextWeekHours.Hours != 2 {
		t.Errorf("Expected 2 hours, got %f", nextWeekHours.Hours)
	}
}
