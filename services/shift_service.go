package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"shift_calc/models"
	"shift_calc/utils"
	"sort"
)

// Read JSON shift data and returns a slice of Shifts
func LoadShifts(filename string) ([]models.Shift, error) {
	// Read the file, make sure it exists and you have permission to read it!
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Make sure the file is not empty!
	if len(data) == 0 {
		return nil, errors.New("this file is empty")
	}

	// Parse the JSON
	var shifts []models.Shift
	err = json.Unmarshal(data, &shifts)
	if err != nil {
		return nil, err
	}

	// There is probably a better way to localize the time and
	// and validate the time format, and required fields
	for i := range shifts {
		// Validate each each shift row has the required fields, and proper time format
		if err := shifts[i].ValidateShift(); err != nil {
			var errorString = fmt.Sprintf("invalid shift data: %v - error: %s", shifts[i], err)
			return nil, errors.New(errorString)
		}

		// Localize the time to CST, so DST is not an issue
		if time, err := utils.LocalizeTime(shifts[i].RawStartTime); err != nil {
			var errorString = fmt.Sprintf("could not localizee starttime in shiftid %d: %s", shifts[i].ShiftID, err)
			return nil, errors.New(errorString)
		} else {
			shifts[i].StartTime = time
		}
		if time, err := utils.LocalizeTime(shifts[i].RawEndTime); err != nil {
			var errorString = fmt.Sprintf("could not localize endtime in shiftid %d: %s", shifts[i].ShiftID, err)
			return nil, errors.New(errorString)
		} else {
			shifts[i].EndTime = time
		}
	}

	return shifts, nil
}

// Condense the shift data in a EmployeeReport
// Requirements:
// Groupd shifts worked by week and employee
// Calculate the total number of regular hours and overtime (>40) worked
// If a shift crosses into the next week split the hours aT midnight and add
// them to the next week summary
// If a shift overlaps with another shift, mark it Invalid and don't count the hours
func CondenseShiftData(shifts []models.Shift) []models.EmployeeReport {
	employeeReport := make(map[int64]map[string]*models.EmployeeReport)

	for _, shift := range shifts {
		// Get the start of the week for the shift, this is one of the keys YYYY-MM-DD
		startOfWeekKey := utils.StartOfWeek(shift.StartTime)

		// Does this exist aready? Create or get
		if _, exists := employeeReport[shift.EmployeeID]; !exists {
			employeeReport[shift.EmployeeID] = make(map[string]*models.EmployeeReport)
		}
		employeeWeek, exists := employeeReport[shift.EmployeeID][startOfWeekKey]
		if !exists {
			employeeWeek = &models.EmployeeReport{
				EmployeeID:    shift.EmployeeID,
				StartOfWeek:   startOfWeekKey,
				RegularHours:  0,
				OverTimeHours: 0,
				InvalidShifts: []int64{},
			}
			employeeReport[shift.EmployeeID][startOfWeekKey] = employeeWeek
		}

		// Get all the past shifts for this employee, loop through them
		// and check for overlaps.
		// Probably a much smarter way to do this given more time... TODO
		var pastShifts []models.Shift
		for _, s := range shifts {
			if s.EmployeeID == shift.EmployeeID && s.ShiftID != shift.ShiftID && utils.StartOfWeek(s.StartTime) == startOfWeekKey {
				pastShifts = append(pastShifts, s)
			}
		}
		if isShiftOverlapping(pastShifts, shift) {
			employeeWeek.InvalidShifts = append(employeeWeek.InvalidShifts, shift.ShiftID)
			continue
		}

		// Now we need to see if the hours cross the midnight mark of the end of the week
		// Return hours for this week start --> midnight and next midnight --> end
		shiftHours, nextWeekShiftHours := utils.ReturnHoursForShift(shift.StartTime, shift.EndTime)

		// Determine if we are in overtime or not
		// Find out how many are left before reaching 40
		// Add to regular hours if we are still under 40
		// Add to overtime if we are over 40
		remainingRegular := 40 - employeeWeek.RegularHours
		if shiftHours.Hours <= remainingRegular {
			employeeWeek.RegularHours += shiftHours.Hours
		} else {
			employeeWeek.RegularHours += remainingRegular
			employeeWeek.OverTimeHours += (shiftHours.Hours - remainingRegular)
		}

		// This code is kind of trash, I should have made a function to handle
		// checking and creating the next week record since it is repeated above TODO
		if nextWeekShiftHours != nil {
			nextWeekStartOfWeek := utils.StartOfWeek(nextWeekShiftHours.Start)

			var nextWeekEmployeeWeek *models.EmployeeReport

			// Does this exist aready? Create or get
			if existingRecord, exists := employeeReport[shift.EmployeeID][nextWeekStartOfWeek]; exists {
				nextWeekEmployeeWeek = existingRecord
			} else {
				nextWeekEmployeeWeek = &models.EmployeeReport{
					EmployeeID:    shift.EmployeeID,
					StartOfWeek:   nextWeekStartOfWeek,
					RegularHours:  0,
					OverTimeHours: 0,
					InvalidShifts: []int64{},
				}
				employeeReport[shift.EmployeeID][nextWeekStartOfWeek] = nextWeekEmployeeWeek
			}

			// Once again repeated code, should have made a function to handle this TODO
			nextWeekRemainingRegular := 40 - nextWeekEmployeeWeek.RegularHours
			if nextWeekShiftHours.Hours <= nextWeekRemainingRegular {
				nextWeekEmployeeWeek.RegularHours += nextWeekShiftHours.Hours
			} else {
				nextWeekEmployeeWeek.RegularHours += nextWeekRemainingRegular
				nextWeekEmployeeWeek.OverTimeHours += (nextWeekShiftHours.Hours - nextWeekRemainingRegular)
			}
		}
	}

	// Flatten the map into a slice and sort by employeeID
	// Don't double loop, sort and append TODO!
	var results []models.EmployeeReport
	for _, weeklySummaries := range employeeReport {
		for _, summary := range weeklySummaries {
			results = append(results, *summary)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].EmployeeID < results[j].EmployeeID
	})

	return results
}

// Check if overlap exists
// If current shifts StartTime is Before the otherShifts EndTime and
// the current shifts EndTime is After the otherShifts StartTime
// This still allows for a shift to be worked (and completed) the same day before the start
// of the other shift, or after the end of the other shift
func isShiftOverlapping(pastShifts []models.Shift, shift models.Shift) bool {
	for _, pastShift := range pastShifts {
		if shift.StartTime.Before(pastShift.EndTime) && shift.EndTime.After(pastShift.StartTime) {
			return true
		}
	}
	return false
}
