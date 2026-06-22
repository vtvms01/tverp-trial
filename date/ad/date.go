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

// ResolveAirDatetime determines air datetime from available sources (DD8 shared helper).
// Priority: planned_start_at > broadcast_start_time + scheduled_date > error.
func ResolveAirDatetime(plannedStartAt *time.Time, broadcastStartTime *string, scheduledDate time.Time) (*time.Time, error) {
	if plannedStartAt != nil {
		t := plannedStartAt.In(vnLocation)
		return &t, nil
	}

	if broadcastStartTime != nil && *broadcastStartTime != "" {
		// Parse TIME "HH:MM" or "0000-01-01THH:MM:SSZ" format from channel_cache
		var hour, min int
		parsed := false
		// Try "HH:MM" first
		if _, err := fmt.Sscanf(*broadcastStartTime, "%d:%d", &hour, &min); err == nil {
			parsed = true
		}
		// Try extracting from timestamp format "0000-01-01T14:00:00Z"
		if !parsed {
			t, err := time.Parse(time.RFC3339, *broadcastStartTime)
			if err == nil {
				hour = t.Hour()
				min = t.Minute()
				parsed = true
			}
		}
		if parsed {
			vnDate := scheduledDate.In(vnLocation)
			combined := time.Date(vnDate.Year(), vnDate.Month(), vnDate.Day(), hour, min, 0, 0, vnLocation)
			return &combined, nil
		}
	}

	return nil, fmt.Errorf("cannot resolve air datetime: no planned_start_at and no broadcast_start_time")
}

// FormatDateVN formats time.Time to "YYYY-MM-DD" string.
func FormatDateVN(t time.Time) string {
	return t.In(vnLocation).Format("2006-01-02")
}
