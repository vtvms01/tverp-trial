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
		vnLocation = time.FixedZone("ICT", 7*60*60)
	}
}

// VNLocation returns the Vietnam timezone location.
func VNLocation() *time.Location {
	return vnLocation
}

// ParseDateInVN parses "YYYY-MM-DD" string and returns time in VN timezone.
func ParseDateInVN(dateStr string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", dateStr, vnLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format '%s': expected YYYY-MM-DD", dateStr)
	}
	return t, nil
}

// NormalizeDateToVN converts any time.Time to VN date-only (00:00:00 VN).
func NormalizeDateToVN(t time.Time) time.Time {
	vnTime := t.In(vnLocation)
	return time.Date(vnTime.Year(), vnTime.Month(), vnTime.Day(), 0, 0, 0, 0, vnLocation)
}

// TodayInVN returns today's date at 00:00:00 in Vietnam timezone.
func TodayInVN() time.Time {
	now := time.Now().In(vnLocation)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, vnLocation)
}

// FormatDateVN formats time.Time to "YYYY-MM-DD" string in VN timezone.
func FormatDateVN(t time.Time) string {
	return t.In(vnLocation).Format("2006-01-02")
}
