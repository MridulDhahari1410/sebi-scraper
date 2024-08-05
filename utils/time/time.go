package time

import (
	"time"

	"sebi-scrapper/constants"
)

var il *time.Location

// Init is used to initialise the time.
func Init() (err error) {
	il, err = time.LoadLocation(constants.IndianTimeLocation)
	return
}

// GetIndianTimeLocation is used to get Indian time location.
func GetIndianTimeLocation() *time.Location {
	return il
}

// GetCurrentTime is used to get the current time.
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

// GetCurrentIndianTime is used to get the current Indian time.
func GetCurrentIndianTime() time.Time {
	return time.Now().In(il)
}

// GetIndianTimeFromEpoch is used to get the time from epoch time in millis.
func GetIndianTimeFromEpoch(t int64) time.Time {
	return time.UnixMilli(t).In(il)
}

// GetPreviousIndianDate is used to get the previous Indian Date.
func GetPreviousIndianDate() time.Time {
	return time.Now().In(il).Add(-time.Hour * 24)
}

// Parse is used to parse time.
func Parse(ts, format string) (time.Time, error) {
	return time.Parse(format, ts)
}

// GetFirstDayOfMonth is used to get the first day of the month.
func GetFirstDayOfMonth() time.Time {
	ct := GetCurrentIndianTime()
	for {
		if ct.Day() == 1 {
			return ct
		}
		ct = ct.Add(-time.Hour * 24)
	}
}

// GetCurrentDateTimeStamp is used to get the Timestamp of th current date as per IST it sets the time to 0.
func GetCurrentDateTimeStamp() time.Time {
	now := GetCurrentIndianTime()

	// Set the time portion to zero
	dateTimeStamp := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, GetIndianTimeLocation())
	return dateTimeStamp
}

// GetFirstDayOfMonthWithZeroTimestamp is used to get the first day of the month with zero timestamp.
func GetFirstDayOfMonthWithZeroTimestamp() time.Time {
	t := GetFirstDayOfMonth()
	ct := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, GetIndianTimeLocation())
	return ct
}
