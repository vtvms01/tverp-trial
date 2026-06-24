// Package util hợp nhất các hàm tiện ích ngày/giờ theo múi giờ Việt Nam
// (Asia/Ho_Chi_Minh) từ 8 biến thể trong date/ (ad, channel, config, editorial,
// nrcs, production, resource, royalty). Xem date/MERGE_NOTES.md để biết cơ sở gộp.
package util

import (
	"fmt"
	"time"
)

// DefaultTimezone là múi giờ chuẩn dùng cho toàn bộ hàm trong package.
const DefaultTimezone = "Asia/Ho_Chi_Minh"

var vnLocation *time.Location

func init() {
	var err error
	vnLocation, err = time.LoadLocation(DefaultTimezone)
	if err != nil {
		vnLocation = time.FixedZone("ICT", 7*60*60)
	}
}

// ───────────────────────── Lõi chung (giống hệt ở cả 8 biến thể) ─────────────────────────

// VNLocation trả về *time.Location của múi giờ Việt Nam.
func VNLocation() *time.Location {
	return vnLocation
}

// ParseDateInVN parse chuỗi "YYYY-MM-DD" và trả về time.Time theo múi giờ VN.
func ParseDateInVN(dateStr string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", dateStr, vnLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format '%s': expected YYYY-MM-DD", dateStr)
	}
	return t, nil
}

// NormalizeDateToVN đưa mọi time.Time về đầu ngày (00:00:00) theo múi giờ VN.
func NormalizeDateToVN(t time.Time) time.Time {
	vnTime := t.In(vnLocation)
	return time.Date(vnTime.Year(), vnTime.Month(), vnTime.Day(), 0, 0, 0, 0, vnLocation)
}

// TodayInVN trả về hôm nay lúc 00:00:00 theo múi giờ VN.
func TodayInVN() time.Time {
	now := time.Now().In(vnLocation)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, vnLocation)
}

// NowInVN trả về thời điểm hiện tại theo múi giờ VN.
func NowInVN() time.Time {
	return time.Now().In(vnLocation)
}

// FormatDateVN format time.Time thành chuỗi "YYYY-MM-DD" theo múi giờ VN.
func FormatDateVN(t time.Time) string {
	return t.In(vnLocation).Format("2006-01-02")
}

// ───────────────────────── So sánh ngày (từ channel) ─────────────────────────

// IsDateInPast cho biết date (chuẩn hóa về ngày VN) có trước hôm nay (VN) hay không.
func IsDateInPast(date time.Time) bool {
	today := TodayInVN()
	normalizedDate := NormalizeDateToVN(date)
	return normalizedDate.Before(today)
}

// IsDateToday cho biết date (chuẩn hóa về ngày VN) có đúng là hôm nay (VN) hay không.
func IsDateToday(date time.Time) bool {
	today := TodayInVN()
	normalizedDate := NormalizeDateToVN(date)
	return normalizedDate.Equal(today)
}

// SameDate cho biết t1 và t2 có rơi vào cùng một ngày dương lịch (VN) hay không.
func SameDate(t1, t2 time.Time) bool {
	d1 := NormalizeDateToVN(t1)
	d2 := NormalizeDateToVN(t2)
	return d1.Equal(d2)
}

// ───────────────────────── Khoảng ngày (từ nrcs) ─────────────────────────

// DateRangeInVN parse [fromStr, toStr] và trả về khoảng nửa mở
// [fromAt, toAtExclusive) trong đó toAtExclusive là đầu ngày kế tiếp của toStr.
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

// ───────────────────────── Format thời gian RFC3339 (từ resource/royalty) ─────────────────────────

// FormatTimeVN format t thành RFC3339 theo múi giờ VN; trả về "" nếu t là zero.
func FormatTimeVN(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(vnLocation).Format(time.RFC3339)
}

// FormatTimeVNPtr là biến thể nil-safe của FormatTimeVN.
func FormatTimeVNPtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatTimeVN(*t)
}

// ───────────────────────── Air datetime (từ ad/editorial/production — 3 bản giống hệt) ─────────────────────────

// ResolveAirDatetime xác định thời điểm phát sóng từ các nguồn sẵn có.
// Ưu tiên: plannedStartAt > broadcastStartTime + scheduledDate > lỗi.
func ResolveAirDatetime(plannedStartAt *time.Time, broadcastStartTime *string, scheduledDate time.Time) (*time.Time, error) {
	if plannedStartAt != nil {
		t := plannedStartAt.In(vnLocation)
		return &t, nil
	}

	if broadcastStartTime != nil && *broadcastStartTime != "" {
		var hour, min int
		parsed := false
		if _, err := fmt.Sscanf(*broadcastStartTime, "%d:%d", &hour, &min); err == nil {
			parsed = true
		}
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
