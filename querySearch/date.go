// These functions are useful when you need to create specific date and time queries,
// like when you need to filter results between a time interval or within a specific time of day.
package querysearch

import (
	"fmt"
	"time"
)

var defaultLocationTimezone = "America/Sao_Paulo"

// Set default timezone for reference
func SetDefaultTimezone(locationDefault string) {
	defaultLocationTimezone = locationDefault
}

// Return date in format '2006-01-02 00:00:00 '
func GetDateFormatted(dateRef time.Time) string {
	dateF := setTimeZone(dateRef)
	return dateF.Format("2006-01-02") + " 00:00:00 " + GetTimeZoneOffSet(dateF)
}

// Return hour and minute in UTC, with offset from the local timezone
func GetTimeZoneOffSet(dateRef time.Time) string {
	_, offset := dateRef.Zone()
	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60
	offsetStr := fmt.Sprintf("%03d:%02d", offsetHours, offsetMinutes)
	return offsetStr
}

// Change date to timezone defined on variable default LocationTimezone
func setTimeZone(dateRef time.Time) time.Time {
	loc, _ := time.LoadLocation(defaultLocationTimezone)
	return dateRef.In(loc)
}

// Define fisrt hour  of day
func SetDateFirstHour(dateRef time.Time) time.Time {
	newDate := setTimeZone(dateRef)
	newDate = time.Date(newDate.Year(), newDate.Month(), newDate.Day(), 0, 0, 0, 0, newDate.Location())
	return newDate
}

// Define last hour of day
func SetDateLastHour(dateRef time.Time) time.Time {
	newDate := setTimeZone(dateRef)
	newDate = time.Date(newDate.Year(), newDate.Month(), newDate.Day(), 23, 59, 59, 0, newDate.Location())
	return newDate
}

// Define the selection of a date >= the date entered
func StartDate(dateRef time.Time) map[string]interface{} {
	return map[string]interface{}{
		"$gte": SetDateFirstHour(dateRef),
	}
}

// Define the selection of a date <= the date entered
func EndDate(dateRef time.Time) map[string]interface{} {
	return map[string]interface{}{
		"$lte": SetDateLastHour(dateRef),
	}
}

// Define the date selection within a period of 1 day from the date entered
func OneDay(dateRef time.Time) map[string]interface{} {
	return map[string]interface{}{
		"$gte": SetDateFirstHour(dateRef),
		"$lte": SetDateLastHour(dateRef),
	}
}

// Define the date range between two dates
func PeriodDay(dateStart, dateEnd time.Time) map[string]interface{} {
	return map[string]interface{}{
		"$gte": SetDateFirstHour(dateStart),
		"$lte": SetDateLastHour(dateEnd),
	}
}
