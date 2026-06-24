package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVNLocation(t *testing.T) {
	loc := VNLocation()
	require.NotNil(t, loc)

	_, offset := time.Now().In(loc).Zone()
	assert.Equal(t, 7*60*60, offset)
}

func TestParseDateInVN(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		year        int
		month       time.Month
		day         int
	}{
		{"ngay hop le", "2026-01-21", false, 2026, time.January, 21},
		{"nam nhuan 29/02", "2024-02-29", false, 2024, time.February, 29},
		{"29/02 nam thuong", "2025-02-29", true, 0, 0, 0},
		{"thang 13", "2026-13-21", true, 0, 0, 0},
		{"ngay 32", "2026-1-32", true, 0, 0, 0},
		{"sai dau phan cach", "2026/01/21", true, 0, 0, 0},
		{"sai thu tu", "21-01-2026", true, 0, 0, 0},
		{"chuoi rong", "", true, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDateInVN(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Kiểm tra ngày tháng
			assert.Equal(t, tt.year, result.Year())
			assert.Equal(t, tt.month, result.Month())
			assert.Equal(t, tt.day, result.Day())

			// Kiểm tra giờ là 00:00:00
			assert.Equal(t, 0, result.Hour())
			assert.Equal(t, 0, result.Minute())
			assert.Equal(t, 0, result.Second())
			assert.Equal(t, 0, result.Nanosecond())

			// Kiểm tra múi giờ Việt Nam
			assert.Equal(t, vnLocation, result.Location(),
				"ParseDateInVN phải trả về time.Time có Location là vnLocation")

			_, offset := result.Zone()
			assert.Equal(t, 7*60*60, offset, "Offset phải là +7 giờ (UTC+7)")
		})
	}
}

func TestNormalizeDateToVN(t *testing.T) {
	vnLoc := VNLocation()
	tests := []struct {
		name      string
		input     time.Time
		wantYear  int
		wantMonth time.Month
		wantDay   int
	}{
		{
			"UTC nua dem -> cung ngay o VN",
			time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC),
			2026, time.January, 21,
		},
		{
			"UTC 23:00 -> sang ngay hom sau o VN",
			time.Date(2026, 1, 20, 23, 0, 0, 0, time.UTC), // 23:00 UTC = 06:00 VN ngay 21
			2026, time.January, 21,
		},
		{
			"gio VN giu nguyen ngay",
			time.Date(2026, 1, 21, 15, 30, 0, 0, vnLoc),
			2026, time.January, 21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeDateToVN(tt.input)
			assert.Equal(t, tt.wantYear, got.Year())
			assert.Equal(t, tt.wantMonth, got.Month())
			assert.Equal(t, tt.wantDay, got.Day())
			assert.Equal(t, 0, got.Hour(), "phai la nua dem")
			_, offset := got.Zone()
			assert.Equal(t, 7*60*60, offset)

		})
	}
}

func TestTodayInVN(t *testing.T) {
	today := TodayInVN()
	assert.Equal(t, 0, today.Hour())
	assert.Equal(t, 0, today.Minute())
	assert.Equal(t, 0, today.Second())
	assert.Equal(t, 0, today.Nanosecond())

	_, offset := today.Zone()
	assert.Equal(t, 7*60*60, offset)
}
func TestFormatDateVN(t *testing.T) {
	// round-trip: parse roi format lai phai ra chuoi goc
	parsed, err := ParseDateInVN("2026-01-21")
	require.NoError(t, err)
	assert.Equal(t, "2026-01-21", FormatDateVN(parsed))

	// thoi diem UTC sat nua dem phai format theo NGAY cua gio VN
	utc := time.Date(2026, 1, 20, 23, 0, 0, 0, time.UTC) // = 21/01 o VN
	assert.Equal(t, "2026-01-21", FormatDateVN(utc))
}

func TestNowInVN(t *testing.T) {
	now := NowInVN()
	_, offset := now.Zone()
	assert.Equal(t, 7*60*60, offset)
}

func TestFormatTimeVN(t *testing.T) {
	assert.Equal(t, "", FormatTimeVN(time.Time{}))

	utc := time.Date(2026, 1, 21, 14, 0, 0, 0, time.UTC)
	assert.Equal(t, "2026-01-21T21:00:00+07:00", FormatTimeVN(utc))
}
func TestFormatTimeVNPtr(t *testing.T) {
	assert.Equal(t, "", FormatTimeVNPtr(nil))

	zero := time.Time{}
	assert.Equal(t, "", FormatTimeVNPtr(&zero))

	utc := time.Date(2026, 1, 21, 14, 0, 0, 0, time.UTC)
	assert.Equal(t, "2026-01-21T21:00:00+07:00", FormatTimeVNPtr(&utc))
}
