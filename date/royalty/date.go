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

// FormatDateVN format time.Time thành chuỗi "YYYY-MM-DD".
func FormatDateVN(t time.Time) string {
	return t.In(vnLocation).Format("2006-01-02")
}

// FormatTimeVN format time.Time thành chuỗi RFC3339 với offset VN +07:00.
// QUAN TRỌNG: Go time.Time.Format() dùng location CỦA CHÍNH t. DB TIMESTAMPTZ khi scan trả về time.Time
// với location UTC sau khi pgx driver giải mã → nếu không .In(vnLocation) tường minh TRƯỚC Format(),
// output sẽ có hậu tố "Z" thay vì "+07:00". Dùng helper này cho trường chuỗi payload outbox + mọi
// serialization timestamp không phải DTO.
//
// Với trường time.Time trong DTO Response (chốt D41 — giữ nguyên time.Time gốc), dùng
// `.In(util.VNLocation())` tại lúc map để Go JSON encoder lấy location VN cho "+07:00".
//
// time.Time rỗng → chuỗi rỗng (bên gọi phân biệt được "chưa đặt" với thời gian thật).
func FormatTimeVN(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(vnLocation).Format(time.RFC3339)
}

// FormatTimeVNPtr là biến thể nullable của FormatTimeVN cho trường *time.Time.
// nil HOẶC zero → chuỗi rỗng.
func FormatTimeVNPtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatTimeVN(*t)
}
