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
			// Dùng := ở đây là đúng vì result và err là biến mới trong scope này
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

func TestResolveAirDateTime(t *testing.T) {
	scheduled := time.Date(2026, 1, 21, 0, 0, 0, 0, VNLocation())
	// NHÁNH 1: có plannedStartAt -> ưu tiên dùng nó

	t.Run("Có plannedStartAt thì dùng nó", func(t *testing.T) {
		planned := time.Date(2026, 1, 21, 14, 30, 0, 0, time.UTC)
		got, err := ResolveAirDatetime(&planned, nil, scheduled)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.True(t, planned.Equal(*got))
		_, offset := got.Zone()
		assert.Equal(t, 7*60*60, offset)

	})
	// NHÁNH 2a: không có planned, có broadcast dạng "HH:MM"
	t.Run("broadcastStartTime dang HH:MM", func(t *testing.T) {
		bt := "14:00"
		got, err := ResolveAirDatetime(nil, &bt, scheduled)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, 14, got.Hour())
		assert.Equal(t, 0, got.Minute())
		assert.Equal(t, 21, got.Day())
	})
	// NHÁNH 3a: broadcast rỗng -> lỗi
	t.Run("broadcastStartTime rong -> loi", func(t *testing.T) {
		empty := ""
		_, err := ResolveAirDatetime(nil, &empty, scheduled)
		assert.Error(t, err)
	})

	// NHÁNH 3b: không có nguồn nào -> lỗi
	t.Run("khong co nguon nao -> loi", func(t *testing.T) {
		_, err := ResolveAirDatetime(nil, nil, scheduled)
		assert.Error(t, err)
	})
	t.Run("broadcastStartTime dang RFC3339", func(t *testing.T) {
		bt := "0000-01-01T09:15:00Z"
		got, err := ResolveAirDatetime(nil, &bt, scheduled)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, 9, got.Hour())
		assert.Equal(t, 15, got.Minute())
	})

}
