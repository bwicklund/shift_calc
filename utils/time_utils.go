package utils

import (
	"shift_calc/models"
	"time"
)

// LocalizeTime makes sure the time is in CST/CDT timezone
// I think this will auto-magically handle the DST issues in the time math...
// Something to check out in testing, the dataset only contains August data
// https://www.reddit.com/r/golang/comments/e6feuc/how_do_you_all_deal_with_daylight_savings_time/
func LocalizeTime(inputTime string) (time.Time, error) {
	centralTimeZone, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return time.Time{}, err
	}

	// Parse the input time to time.Time and RFC3339 format
	newTime, err := time.Parse(time.RFC3339, inputTime)
	if err != nil {
		return time.Time{}, err
	}

	localTime := newTime.In(centralTimeZone)
	return localTime, nil
}

func StartOfWeek(inputTime time.Time) string {
	// Set to CST/CDT timezone
	loc, _ := time.LoadLocation("America/Chicago")
	inputTime = inputTime.In(loc)

	// Set the time to start of the week Sunday at midnight
	offset := -int(inputTime.Weekday())
	startOfWeek := inputTime.AddDate(0, 0, offset)

	// Return formatted date string (YYYY-MM-DD)
	return startOfWeek.Format("2006-01-02") // Go's reference date formatting
}

// Same as above but returns time.Time for math
// I should fix this up and make it a single function
// that returns a time.Time and allows the consumer to format it
func StartOfWeekTime(inputTime time.Time) time.Time {
	// Set to CST/CDT timezone
	loc, _ := time.LoadLocation("America/Chicago")
	inputTime = inputTime.In(loc)

	// Set the time to start of the week Sunday at midnight
	offset := -int(inputTime.Weekday())
	startOfWeek := inputTime.AddDate(0, 0, offset)

	// Return time.Time
	return time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, loc)
}

// Given a shift start and end time split the shift at midnight
// I think this will auto-magically handle the DST issues in the time math...
// Return this weeks hours and next weeks hours if the shift does cross the week boundary
func ReturnHoursForShift(start, end time.Time) (models.Hours, *models.Hours) {
	var currentWeekHours models.Hours
	var nextWeekHours *models.Hours = nil

	midnightThisWeek := StartOfWeekTime(start)
	midnightNextWeek := midnightThisWeek.Add(7 * 24 * time.Hour)

	// This is if the shift is contained within the same week
	if end.Before(midnightNextWeek) || end.Equal(midnightNextWeek) {
		currentWeekHours = models.Hours{
			Start: start,
			End:   end,
			Hours: end.Sub(start).Hours(),
		}
		return currentWeekHours, nil
	}

	// If the shift crosses that boundry we return both sets of hours
	currentWeekHours = models.Hours{
		Start: start,
		End:   midnightNextWeek,
		Hours: midnightNextWeek.Sub(start).Hours(),
	}

	nextWeekHours = &models.Hours{
		Start: midnightNextWeek,
		End:   end,
		Hours: end.Sub(midnightNextWeek).Hours(),
	}

	return currentWeekHours, nextWeekHours
}
