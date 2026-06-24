package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ───────────────────────── Lõi chung ─────────────────────────

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

			assert.Equal(t, tt.year, result.Year())
			assert.Equal(t, tt.month, result.Month())
			assert.Equal(t, tt.day, result.Day())

			assert.Equal(t, 0, result.Hour())
			assert.Equal(t, 0, result.Minute())
			assert.Equal(t, 0, result.Second())
			assert.Equal(t, 0, result.Nanosecond())

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

func TestNowInVN(t *testing.T) {
	now := NowInVN()
	_, offset := now.Zone()
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

// ───────────────────────── So sánh ngày (channel) ─────────────────────────

func TestIsDateInPast(t *testing.T) {
	today := TodayInVN()
	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"hom qua la qua khu", today.AddDate(0, 0, -1), true},
		{"hom nay khong phai qua khu", today, false},
		{"ngay mai khong phai qua khu", today.AddDate(0, 0, 1), false},
		{"tuan truoc la qua khu", today.AddDate(0, 0, -7), true},
		{"tuan sau khong phai qua khu", today.AddDate(0, 0, 7), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsDateInPast(tt.date))
		})
	}
}

func TestIsDateToday(t *testing.T) {
	vnLoc := VNLocation()
	today := TodayInVN()

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"hom nay (dau ngay)", today, true},
		{"hom qua", today.AddDate(0, 0, -1), false},
		{"ngay mai", today.AddDate(0, 0, 1), false},
		{
			name:     "hom nay nhung co gio phut van true",
			date:     time.Date(today.Year(), today.Month(), today.Day(), 15, 30, 0, 0, vnLoc),
			expected: true,
		},
		{
			// UTC 00:00 cung ngay-VN: normalize ve VN van ra hom nay
			name:     "UTC 00:00 cua ngay hom nay (VN)",
			date:     time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC),
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsDateToday(tt.date))
		})
	}
}

func TestSameDate(t *testing.T) {
	vnLoc := VNLocation()

	tests := []struct {
		name     string
		t1, t2   time.Time
		expected bool
	}{
		{
			name:     "cung ngay, khac gio -> true",
			t1:       time.Date(2026, 1, 21, 8, 0, 0, 0, vnLoc),
			t2:       time.Date(2026, 1, 21, 23, 59, 59, 0, vnLoc),
			expected: true,
		},
		{
			name:     "khac ngay -> false",
			t1:       time.Date(2026, 1, 21, 12, 0, 0, 0, vnLoc),
			t2:       time.Date(2026, 1, 22, 12, 0, 0, 0, vnLoc),
			expected: false,
		},
		{
			name:     "cung thoi diem chinh xac -> true",
			t1:       time.Date(2026, 1, 21, 10, 30, 0, 0, vnLoc),
			t2:       time.Date(2026, 1, 21, 10, 30, 0, 0, vnLoc),
			expected: true,
		},
		{
			// 2026-01-20 23:00 UTC = 2026-01-21 06:00 VN -> cung ngay VN
			name:     "lech mui gio nhung cung NGAY VN -> true",
			t1:       time.Date(2026, 1, 20, 23, 0, 0, 0, time.UTC),
			t2:       time.Date(2026, 1, 21, 9, 0, 0, 0, vnLoc),
			expected: true,
		},
		{
			// 2026-01-21 18:00 UTC = 2026-01-22 01:00 VN -> sang ngay khac o VN
			name:     "cung ngay UTC nhung KHAC ngay VN -> false",
			t1:       time.Date(2026, 1, 21, 18, 0, 0, 0, time.UTC),
			t2:       time.Date(2026, 1, 21, 9, 0, 0, 0, vnLoc),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SameDate(tt.t1, tt.t2))
		})
	}
}

// ───────────────────────── Khoảng ngày (nrcs) ─────────────────────────

func TestDateRangeInVN(t *testing.T) {
	t.Run("Khoang hop le nhieu ngay", func(t *testing.T) {
		from, toExcl, err := DateRangeInVN("2026-01-21", "2026-01-23")
		require.NoError(t, err)
		assert.Equal(t, 21, from.Day())
		assert.Equal(t, 0, from.Hour())
		assert.Equal(t, 24, toExcl.Day())
		assert.Equal(t, 0, toExcl.Hour())
	})
	t.Run("from == to (1 ngay) van hop le", func(t *testing.T) {
		from, toExcl, err := DateRangeInVN("2026-01-21", "2026-01-21")
		require.NoError(t, err)
		assert.Equal(t, 21, from.Day())
		assert.Equal(t, 22, toExcl.Day())
	})
	t.Run("from sai dinh dang => loi", func(t *testing.T) {
		_, _, err := DateRangeInVN("xxx", "2026-01-21")
		assert.Error(t, err)
	})
	t.Run("to sai dinh dang => loi", func(t *testing.T) {
		_, _, err := DateRangeInVN("2026-01-21", "")
		assert.Error(t, err)
	})
}

// ───────────────────────── Format thời gian RFC3339 (resource/royalty) ─────────────────────────

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

// ───────────────────────── Air datetime (ad/editorial/production) ─────────────────────────

func TestResolveAirDatetime(t *testing.T) {
	scheduled := time.Date(2026, 1, 21, 0, 0, 0, 0, VNLocation())

	t.Run("Co plannedStartAt thi dung no", func(t *testing.T) {
		planned := time.Date(2026, 1, 21, 14, 30, 0, 0, time.UTC)
		got, err := ResolveAirDatetime(&planned, nil, scheduled)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.True(t, planned.Equal(*got))
		_, offset := got.Zone()
		assert.Equal(t, 7*60*60, offset)
	})
	t.Run("broadcastStartTime dang HH:MM", func(t *testing.T) {
		bt := "14:00"
		got, err := ResolveAirDatetime(nil, &bt, scheduled)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, 14, got.Hour())
		assert.Equal(t, 0, got.Minute())
		assert.Equal(t, 21, got.Day())
	})
	t.Run("broadcastStartTime dang RFC3339", func(t *testing.T) {
		bt := "0000-01-01T09:15:00Z"
		got, err := ResolveAirDatetime(nil, &bt, scheduled)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, 9, got.Hour())
		assert.Equal(t, 15, got.Minute())
	})
	t.Run("broadcastStartTime rong -> loi", func(t *testing.T) {
		empty := ""
		_, err := ResolveAirDatetime(nil, &empty, scheduled)
		assert.Error(t, err)
	})
	t.Run("khong co nguon nao -> loi", func(t *testing.T) {
		_, err := ResolveAirDatetime(nil, nil, scheduled)
		assert.Error(t, err)
	})
}
