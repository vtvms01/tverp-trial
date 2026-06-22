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

// NowInVN returns current time in Vietnam timezone.
func NowInVN() time.Time {
	return time.Now().In(vnLocation)
}

// FormatDateVN formats time.Time to "YYYY-MM-DD" string.
func FormatDateVN(t time.Time) string {
	return t.In(vnLocation).Format("2006-01-02")
}

// FormatTimeVN formats time.Time to RFC3339 string with +07:00 VN offset.
// CRITICAL: Go time.Time.Format() uses t's OWN location. DB TIMESTAMPTZ scan returns time.Time
// with UTC location after pgx driver decode → without explicit .In(vnLocation) BEFORE Format(),
// output is "Z" suffix not "+07:00". Use this helper for outbox payload string fields + any
// non-DTO timestamp serialization.
//
// For DTO Response time.Time fields (D41 lock — preserve native time.Time), use
// `.In(util.VNLocation())` at mapper time so Go JSON encoder picks up VN location for "+07:00".
//
// Empty time.Time → empty string (caller can distinguish "not set" from real time).
func FormatTimeVN(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(vnLocation).Format(time.RFC3339)
}

// FormatTimeVNPtr is the nullable counterpart of FormatTimeVN for *time.Time fields.
// nil OR zero → empty string.
func FormatTimeVNPtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatTimeVN(*t)
}
