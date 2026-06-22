// Package util provides shared date utilities.
// Date helpers — Asia/Ho_Chi_Minh timezone canonical.
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

// DateRangeInVN parses YYYY-MM-DD from/to strings into a half-open VN-timezone range
// Half-open date range [from, to_exclusive) for range queries:
//
//   - from_date: start of day in Asia/Ho_Chi_Minh — inclusive (00:00:00+07:00)
//   - to_date:   start of NEXT day in Asia/Ho_Chi_Minh — exclusive (next-day 00:00:00+07:00)
//
// Query callers MUST use timestamp range `WHERE event_at >= fromAt AND event_at < toAtExclusive`
// to keep PG idx on TIMESTAMPTZ usable (avoid DATE(col) wrapping).
//
// Returns (zero, zero, err) on parse failure or inverted range (fromAt >= toAtExclusive)
// — F-R1-6 guard prevents silent empty results when user accidentally swaps params.
func DateRangeInVN(fromStr, toStr string) (fromAt, toAtExclusive time.Time, err error) {
	from, err := ParseDateInVN(fromStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("from_date: %w", err)
	}
	to, err := ParseDateInVN(toStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("to_date: %w", err)
	}
	toAtExclusive = to.AddDate(0, 0, 1)
	if !from.Before(toAtExclusive) {
		return time.Time{}, time.Time{}, fmt.Errorf("to_date must be after from_date")
	}
	return from, toAtExclusive, nil
}
