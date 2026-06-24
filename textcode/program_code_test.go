package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeVietnamese(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"chữ D hoa", "Đường", "Duong"},
		{"chữ d thường", "đỉnh", "dinh"},
		{"dấu thanh", "Tiếng việt", "Tieng viet"},
		{"nguyên âm đặc biệt", "ăâêôơư", "aaeoou"},
		{"chuỗi rỗng", "", ""},
		{"không dấu giữ nguyên", "ABC 123", "ABC 123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, normalizeVietnamese(tt.in))
		})
	}

}
func TestGenerateExpectedProgramCode(t *testing.T) {
	tests := []struct{ name, dept, cat, title, want string }{
		{"binh thuong", "VTV", "TT", "Thời sự", "VTV_TT_THOISU"},
		{"chu D", "VTV", "GT", "Đường lên đỉnh", "VTV_GT_DUONGLENDINH"},
		{"title rong co dau _ thua", "VTV", "TT", "", "VTV_TT_"},
		{"nhieu space bi xoa", "VTV", "TT", "A   B", "VTV_TT_AB"},
		{"ky tu dac biet bi xoa", "VTV", "TT", "Tin!!! @nóng", "VTV_TT_TINNONG"},
		// phản-ví-dụ: dept thường bị xóa sạch — ghi nhận hành vi thật:
		{"dept thuong bi xoa", "vtv", "TT", "Tin", "_TT_TIN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateExpectedProgramCode(tt.dept, tt.cat, tt.title)
			assert.Equal(t, tt.want, got)
			assert.LessOrEqual(t, len(got), 30, "khong duoc vuot 30 ky tu")
		})
	}
}

func TestGenerateProgramCode_Deterministic(t *testing.T) {
	a := GenerateExpectedProgramCode("VTV", "TT", "Chương trình ABC")
	b := GenerateExpectedProgramCode("VTV", "TT", "Chương trình ABC")
	assert.Equal(t, a, b)
}

func TestGenerateProgramCode_Truncate30(t *testing.T) {
	long := "Day la mot tieu de cuc ky dai vuot qua ba muoi ky tu"
	got := GenerateExpectedProgramCode("VTV", "TT", long)
	assert.Len(t, got, 30)
}
