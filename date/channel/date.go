// FILE: pkg/util/date.go
// Các hàm tiện ích múi giờ cho Việt Nam (Asia/Ho_Chi_Minh)

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
		// Dự phòng về UTC+7
		vnLocation = time.FixedZone("ICT", 7*60*60)
	}
}

// VNLocation trả về *time.Location của múi giờ Việt Nam
func VNLocation() *time.Location {
	return vnLocation
}

// TodayInVN trả về ngày hôm nay lúc 00:00:00 theo múi giờ Việt Nam
func TodayInVN() time.Time {
	now := time.Now().In(vnLocation)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, vnLocation)
}

// ParseDateInVN parse chuỗi "YYYY-MM-DD" và trả về time theo múi giờ VN
func ParseDateInVN(dateStr string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", dateStr, vnLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format '%s': expected YYYY-MM-DD", dateStr)
	}
	return t, nil
}

// NormalizeDateToVN đưa mọi time.Time về chỉ-ngày theo VN (00:00:00 VN)
// QUAN TRỌNG: Dùng hàm này khi đọc DATE từ PostgreSQL để tránh lệch UTC
// Trước tiên đổi sang múi giờ VN, sau đó trích các thành phần ngày
func NormalizeDateToVN(t time.Time) time.Time {
	vnTime := t.In(vnLocation)
	return time.Date(vnTime.Year(), vnTime.Month(), vnTime.Day(), 0, 0, 0, 0, vnLocation)
}

// IsDateInPast kiểm tra date có trước hôm nay không (múi giờ VN)
func IsDateInPast(date time.Time) bool {
	today := TodayInVN()
	// Chuẩn hóa cả hai về VN để so sánh chính xác
	normalizedDate := NormalizeDateToVN(date)
	return normalizedDate.Before(today)
}

// IsDateToday kiểm tra date có phải hôm nay không (múi giờ VN)
func IsDateToday(date time.Time) bool {
	today := TodayInVN()
	normalizedDate := NormalizeDateToVN(date)
	return normalizedDate.Equal(today)
}

// FormatDateVN format time.Time thành chuỗi "YYYY-MM-DD"
func FormatDateVN(t time.Time) string {
	return t.In(vnLocation).Format("2006-01-02")
}

// ToLocalDate là alias của NormalizeDateToVN để tương thích ngược
// KHÔNG NÊN DÙNG: Hãy dùng NormalizeDateToVN
var ToLocalDate = NormalizeDateToVN

// SameDate kiểm tra hai mốc thời gian có cùng ngày dương lịch theo múi giờ VN
func SameDate(t1, t2 time.Time) bool {
	d1 := NormalizeDateToVN(t1)
	d2 := NormalizeDateToVN(t2)
	return d1.Equal(d2)
}
