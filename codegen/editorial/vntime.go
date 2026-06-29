package util

import "time"

// VNLocation trả về location Asia/Ho_Chi_Minh (UTC+7), với dự phòng
// fixed-zone ICT nếu tzdata không khả dụng.
func VNLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return time.FixedZone("ICT", 7*60*60)
	}
	return loc
}

// TodayInVN trả về nửa đêm (00:00) của ngày hiện tại theo múi giờ VN.
func TodayInVN() time.Time {
	loc := VNLocation()
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
}
