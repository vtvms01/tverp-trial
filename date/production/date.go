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

// VNLocation trả về *time.Location của múi giờ Việt Nam.
func VNLocation() *time.Location {
	return vnLocation
}

// ParseDateInVN parse chuỗi "YYYY-MM-DD" và trả về time theo múi giờ VN.
func ParseDateInVN(dateStr string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", dateStr, vnLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format '%s': expected YYYY-MM-DD", dateStr)
	}
	return t, nil
}

// NormalizeDateToVN đưa mọi time.Time về chỉ-ngày theo VN (00:00:00 VN).
func NormalizeDateToVN(t time.Time) time.Time {
	vnTime := t.In(vnLocation)
	return time.Date(vnTime.Year(), vnTime.Month(), vnTime.Day(), 0, 0, 0, 0, vnLocation)
}

// TodayInVN trả về ngày hôm nay lúc 00:00:00 theo múi giờ Việt Nam.
func TodayInVN() time.Time {
	now := time.Now().In(vnLocation)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, vnLocation)
}

// NowInVN trả về thời điểm hiện tại theo múi giờ Việt Nam.
func NowInVN() time.Time {
	return time.Now().In(vnLocation)
}

// ResolveAirDatetime xác định thời điểm phát sóng từ các nguồn sẵn có (helper dùng chung DD8).
// Ưu tiên: planned_start_at > broadcast_start_time + scheduled_date > lỗi.
func ResolveAirDatetime(plannedStartAt *time.Time, broadcastStartTime *string, scheduledDate time.Time) (*time.Time, error) {
	if plannedStartAt != nil {
		t := plannedStartAt.In(vnLocation)
		return &t, nil
	}

	if broadcastStartTime != nil && *broadcastStartTime != "" {
		// Parse TIME dạng "HH:MM" hoặc "0000-01-01THH:MM:SSZ" từ channel_cache
		var hour, min int
		parsed := false
		// Thử "HH:MM" trước
		if _, err := fmt.Sscanf(*broadcastStartTime, "%d:%d", &hour, &min); err == nil {
			parsed = true
		}
		// Thử trích từ định dạng timestamp "0000-01-01T14:00:00Z"
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

// FormatDateVN format time.Time thành chuỗi "YYYY-MM-DD".
func FormatDateVN(t time.Time) string {
	return t.In(vnLocation).Format("2006-01-02")
}
