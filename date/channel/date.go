// FILE: pkg/util/date.go
// Timezone utilities for Vietnam (Asia/Ho_Chi_Minh)

package util

import (
	"fmt"
	"time"
)

const DefaultTimezone = "Asia/Ho_Chi_Minh"

var vnLocation *time.Location

func init() {
	var err error
	vnLocation, err = time.LoadLocation(DefaultTimezone)
	if err != nil {
		// Fallback to UTC+7
		vnLocation = time.FixedZone("ICT", 7*60*60)
	}
}

// VNLocation returns the Vietnam timezone location
func VNLocation() *time.Location {
	return vnLocation
}

// TodayInVN returns today's date at 00:00:00 in Vietnam timezone
func TodayInVN() time.Time {
	now := time.Now().In(vnLocation)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, vnLocation)
}

// ParseDateInVN parses "YYYY-MM-DD" string and returns time in VN timezone
func ParseDateInVN(dateStr string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", dateStr, vnLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format '%s': expected YYYY-MM-DD", dateStr)
	}
	return t, nil
}

// NormalizeDateToVN converts any time.Time to VN date-only (00:00:00 VN)
// CRITICAL: Use this when reading DATE from PostgreSQL to avoid UTC offset issues
// First converts to VN timezone, then extracts date components
func NormalizeDateToVN(t time.Time) time.Time {
	vnTime := t.In(vnLocation)
	return time.Date(vnTime.Year(), vnTime.Month(), vnTime.Day(), 0, 0, 0, 0, vnLocation)
}

// IsDateInPast checks if the given date is before today (VN timezone)
func IsDateInPast(date time.Time) bool {
	today := TodayInVN()
	// Normalize both to VN for accurate comparison
	normalizedDate := NormalizeDateToVN(date)
	return normalizedDate.Before(today)
}

// IsDateToday checks if the given date is today (VN timezone)
func IsDateToday(date time.Time) bool {
	today := TodayInVN()
	normalizedDate := NormalizeDateToVN(date)
	return normalizedDate.Equal(today)
}

// FormatDateVN formats time.Time to "YYYY-MM-DD" string
func FormatDateVN(t time.Time) string {
	return t.In(vnLocation).Format("2006-01-02")
}

// ToLocalDate is an alias for NormalizeDateToVN for backward compatibility
// DEPRECATED: Use NormalizeDateToVN instead
var ToLocalDate = NormalizeDateToVN

// SameDate checks if two times represent the same calendar date in VN timezone
func SameDate(t1, t2 time.Time) bool {
	d1 := NormalizeDateToVN(t1)
	d2 := NormalizeDateToVN(t2)
	return d1.Equal(d2)
}
