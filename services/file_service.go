package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"shift_calc/models"
	"shift_calc/utils"
)

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

	var shifts []models.Shift
	err = json.Unmarshal(data, &shifts)
	if err != nil {
		return nil, err
	}

	// There is probably a better way to localize the time and
	// and validate the time format, but this works for now.
	for _, shift := range shifts {

		// Validate each each shift row has the required fields, and proper time format
		if err := shift.ValidateShift(); err != nil {
			var errorString = fmt.Sprintf("Invalid shift data: %v - Error: %s", shift, err)
			return nil, errors.New(errorString)
		}

		// Localize the time to CST, so DST is not an issue
		if time, err := utils.LocalizeTime(shift.StartTime); err != nil {
			var errorString = fmt.Sprintf("could not localizee StartTime in ShiftID %d: %s", shift.ShiftID, err)
			return nil, errors.New(errorString)
		} else {
			shift.StartTime = time
		}
		if time, err := utils.LocalizeTime(shift.EndTime); err != nil {
			var errorString = fmt.Sprintf("could not localize EndTime in ShiftID %d: %s", shift.ShiftID, err)
			return nil, errors.New(errorString)
		} else {
			shift.EndTime = time
		}
	}

	return shifts, nil
}
