package util

import "time"

// VNLocation returns the Asia/Ho_Chi_Minh location (UTC+7), with an ICT
// fixed-zone fallback if tzdata is unavailable.
func VNLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return time.FixedZone("ICT", 7*60*60)
	}
	return loc
}

// TodayInVN returns midnight (00:00) of the current day in VN timezone.
func TodayInVN() time.Time {
	loc := VNLocation()
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
}
