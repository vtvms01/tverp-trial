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

// FormatDateVN format time.Time thành chuỗi "YYYY-MM-DD" theo múi giờ VN.
func FormatDateVN(t time.Time) string {
	return t.In(vnLocation).Format("2006-01-02")
}
