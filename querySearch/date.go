package querysearch

import (
	"fmt"
	"time"
)

var defaultLocationTimezone = "America/Sao_Paulo"

func SetDefaultTimezone(locationDefault string) {
	defaultLocationTimezone = locationDefault
}

func GetDateFormatted(dateRef time.Time) string {
	dateF := setTimeZone(dateRef)
	return dateF.Format("2006-01-02") + " 00:00:00 " + GetTimeZoneOffSet(dateF)
}

func GetTimeZoneOffSet(dateRef time.Time) string {
	_, offset := dateRef.Zone()
	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60
	offsetStr := fmt.Sprintf("%03d:%02d", offsetHours, offsetMinutes)
	return offsetStr
}

func setTimeZone(dateRef time.Time) time.Time {
	loc, _ := time.LoadLocation(defaultLocationTimezone)
	return dateRef.In(loc)
}

func SetDateFirstHour(dateRef time.Time) time.Time {
	newDate := setTimeZone(dateRef)
	newDate = time.Date(newDate.Year(), newDate.Month(), newDate.Day(), 0, 0, 0, 0, newDate.Location())
	return newDate
}

func SetDateLastHour(dateRef time.Time) time.Time {
	newDate := setTimeZone(dateRef)
	newDate = time.Date(newDate.Year(), newDate.Month(), newDate.Day(), 23, 59, 59, 0, newDate.Location())
	return newDate
}

func StartDate(dateRef time.Time) map[string]interface{} {
	return map[string]interface{}{
		"$gte": SetDateFirstHour(dateRef),
	}
}

func EndDate(dateRef time.Time) map[string]interface{} {
	return map[string]interface{}{
		"$lte": SetDateLastHour(dateRef),
	}
}

func OneDay(dateRef time.Time) map[string]interface{} {
	return map[string]interface{}{
		"$gte": SetDateFirstHour(dateRef),
		"$lte": SetDateLastHour(dateRef),
	}
}

func PeriodDay(dateStart, dateEnd time.Time) map[string]interface{} {
	return map[string]interface{}{
		"$gte": SetDateFirstHour(dateStart),
		"$lte": SetDateLastHour(dateEnd),
	}
}

func ConvertToRFC3339(dateRef string) (string, error) {
	dateFront, err := time.Parse(time.RFC3339, dateRef)
	if err != nil {
		return "", err
	}
	return dateFront.UTC().Format(time.RFC3339), nil
}
