package utils

import (
	"time"
)

func SetTimeZone(timeZoneFrom, timeZoneTo time.Time) time.Time {
	year, month, day := timeZoneTo.Date()
	hour, min, sec := timeZoneTo.Clock()
	nsec := timeZoneTo.Nanosecond()

	return time.Date(year, month, day, hour, min, sec, nsec, timeZoneFrom.Location())
}
