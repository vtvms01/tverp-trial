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

	// VN is UTC+7
	now := time.Now().In(loc)
	_, offset := now.Zone()
	assert.Equal(t, 7*60*60, offset, "Vietnam should be UTC+7")
}

func TestTodayInVN(t *testing.T) {
	today := TodayInVN()

	// Should be midnight
	assert.Equal(t, 0, today.Hour())
	assert.Equal(t, 0, today.Minute())
	assert.Equal(t, 0, today.Second())
	assert.Equal(t, 0, today.Nanosecond())

	// Should be in VN timezone
	_, offset := today.Zone()
	assert.Equal(t, 7*60*60, offset)
}

func TestParseDateInVN(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectYear  int
		expectMonth time.Month
		expectDay   int
	}{
		{
			name:        "valid date",
			input:       "2026-01-21",
			expectError: false,
			expectYear:  2026,
			expectMonth: time.January,
			expectDay:   21,
		},
		{
			name:        "another valid date",
			input:       "2025-12-31",
			expectError: false,
			expectYear:  2025,
			expectMonth: time.December,
			expectDay:   31,
		},
		{
			name:        "invalid format - wrong separator",
			input:       "2026/01/21",
			expectError: true,
		},
		{
			name:        "invalid format - wrong order",
			input:       "21-01-2026",
			expectError: true,
		},
		{
			name:        "invalid date",
			input:       "2026-13-01",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDateInVN(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectYear, result.Year())
				assert.Equal(t, tt.expectMonth, result.Month())
				assert.Equal(t, tt.expectDay, result.Day())
				assert.Equal(t, 0, result.Hour())
				assert.Equal(t, 0, result.Minute())
				assert.Equal(t, 0, result.Second())

				// Should be in VN timezone
				_, offset := result.Zone()
				assert.Equal(t, 7*60*60, offset)
			}
		})
	}
}

func TestNormalizeDateToVN(t *testing.T) {
	vnLoc := VNLocation()

	tests := []struct {
		name           string
		input          time.Time
		expectedYear   int
		expectedMonth  time.Month
		expectedDay    int
		expectedOffset int
	}{
		{
			name:           "UTC midnight → VN date",
			input:          time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC),
			expectedYear:   2026,
			expectedMonth:  time.January,
			expectedDay:    21,
			expectedOffset: 7 * 60 * 60,
		},
		{
			name:           "UTC 23:00 → VN next day (correctly converts timezone first)",
			input:          time.Date(2026, 1, 20, 23, 0, 0, 0, time.UTC), // 23:00 UTC = 06:00+1 VN
			expectedYear:   2026,
			expectedMonth:  time.January,
			expectedDay:    21, // 2026-01-20 23:00 UTC = 2026-01-21 06:00 VN → normalizes to 2026-01-21
			expectedOffset: 7 * 60 * 60,
		},
		{
			name:           "VN time preserves date",
			input:          time.Date(2026, 1, 21, 15, 30, 0, 0, vnLoc),
			expectedYear:   2026,
			expectedMonth:  time.January,
			expectedDay:    21,
			expectedOffset: 7 * 60 * 60,
		},
		{
			name:           "VN midnight stays same",
			input:          time.Date(2026, 1, 21, 0, 0, 0, 0, vnLoc),
			expectedYear:   2026,
			expectedMonth:  time.January,
			expectedDay:    21,
			expectedOffset: 7 * 60 * 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeDateToVN(tt.input)

			assert.Equal(t, tt.expectedYear, result.Year())
			assert.Equal(t, tt.expectedMonth, result.Month())
			assert.Equal(t, tt.expectedDay, result.Day())
			assert.Equal(t, 0, result.Hour(), "should be midnight")
			assert.Equal(t, 0, result.Minute())
			assert.Equal(t, 0, result.Second())

			_, offset := result.Zone()
			assert.Equal(t, tt.expectedOffset, offset)
		})
	}
}

func TestIsDateInPast(t *testing.T) {
	today := TodayInVN()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"yesterday is in past", yesterday, true},
		{"today is not in past", today, false},
		{"tomorrow is not in past", tomorrow, false},
		{"last week is in past", today.AddDate(0, 0, -7), true},
		{"next week is not in past", today.AddDate(0, 0, 7), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsDateInPast(tt.date))
		})
	}
}

func TestIsDateInPast_TimezoneEdgeCases(t *testing.T) {
	vnLoc := VNLocation()

	t.Run("UTC date that is today in VN", func(t *testing.T) {
		// Get today in VN
		today := TodayInVN()

		// Same date but in UTC (could be different actual time)
		utcDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

		// Should NOT be in past because we normalize to VN date components
		assert.False(t, IsDateInPast(utcDate))
	})

	t.Run("VN date that is definitely yesterday", func(t *testing.T) {
		today := TodayInVN()
		yesterday := time.Date(today.Year(), today.Month(), today.Day()-1, 23, 59, 59, 0, vnLoc)

		assert.True(t, IsDateInPast(yesterday))
	})
}

func TestDateFormatConsistency(t *testing.T) {
	t.Run("parse and format round-trip", func(t *testing.T) {
		original := "2026-01-21"
		parsed, err := ParseDateInVN(original)
		require.NoError(t, err)

		formatted := parsed.Format("2006-01-02")
		assert.Equal(t, original, formatted)
	})
}
