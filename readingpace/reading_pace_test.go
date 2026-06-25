package timing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandNumber_Basic(t *testing.T) {
	assert.Equal(t, "nam", expandNumber("5"))
}

func TestExpandnumber_TensDigits(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"20", "hai muoi"},
		{"30", "ba muoi"},
		{"60", "sau muoi"},
		{"70", "bay muoi"},
		{"75", "bay muoi lam"},
		{"80", "tam muoi"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, expandNumber(tt.in))
		})
	}
}

func TestPronunciationDict_AdjustedWordCount(t *testing.T) {
	cfg := ReadingPaceConfig{
		DefaultWPM:        160,
		PronunciationDict: map[string]string{"Zelensky": "de len ski"},
	}
	r := CalculatePace("Tong thong Zelensky phat bieu", cfg)

	// Suy luận: câu có 5 từ, 1 từ khó "Zelensky".
	// Từ khó dài 1.2× → đã đếm 1.0, chỉ cộng THÊM 0.2 → 5.2 → làm tròn 5.
	assert.Equal(t, 5, r.WordCount)
	assert.Equal(t, 5, r.AdjustedWordCount) // code bug cho ra 6
}

func TestPronunciationDict_NoDict(t *testing.T) {
	cfg := ReadingPaceConfig{DefaultWPM: 160}
	r := CalculatePace("Tong thong zelensky phat bieu", cfg)
	assert.Equal(t, r.WordCount, r.AdjustedWordCount)
}
