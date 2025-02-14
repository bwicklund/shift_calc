package utils

import (
	"time"
)

// LocalizeTime makes sure the time is in CST/CDT timezone
// I think this will auto-magically handle the DST issues in the time math...
// https://www.reddit.com/r/golang/comments/e6feuc/how_do_you_all_deal_with_daylight_savings_time/
func LocalizeTime(inputTime time.Time) (time.Time, error) {
	centralTimeZone, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return time.Time{}, err
	}
	return inputTime.In(centralTimeZone), nil
}
