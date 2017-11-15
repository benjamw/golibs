package random

import (
	"strconv"
	"time"
)

// NOTE: for format parsing, 1136239445 is the reference unix time
// Mon Jan 2 15:04:05 MST 2006 == 2006-01-02 15:04:05 -0700 == 01/02 3:04:05 pm '06 -0700

// Daten returns a random date between min year and max year
func Daten(min int, max int) time.Time {
	minTime, err := time.Parse("2006-01-02", strconv.Itoa(min)+"-01-01")
	minUnix := minTime.Unix()
	if err != nil {
		minUnix = 0
	}

	maxTime, err := time.Parse("2006-01-02T15:04:05", strconv.Itoa(max)+"-12-31T23:59:59")
	maxUnix := maxTime.Unix()
	if err != nil {
		maxUnix = 2147483647
	}

	return time.Unix(Intn(minUnix, maxUnix), 0)
}

// Date returns a random date between 1999 and 2020
func Date() time.Time {
	return Daten(1999, 2020)
}
