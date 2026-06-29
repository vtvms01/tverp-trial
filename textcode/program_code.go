// FILE: pkg/util/program_code.go
// Package textcode sinh và chuẩn hóa mã nội dung.
//
// Xử lý dấu tiếng Việt: Đ/đ KHÔNG phân rã qua NFD, nên được thay thế tường minh
// TRƯỚC khi chuẩn hóa Unicode.
package util

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	nonAlphaNumericRegex = regexp.MustCompile(`[^A-Z0-9 ]+`)
)

// GenerateExpectedProgramCode sinh mã kỳ vọng cho một chương trình.
// Định dạng: {DEPT}_{CAT}_{TITLE}, tối đa 30 ký tự.
// Tính xác định: cùng title luôn cho ra cùng mã.
func GenerateExpectedProgramCode(departmentCode, categoryCode, title string) string {
	prefix := departmentCode + "_" + categoryCode

	sanitized := normalizeVietnamese(title)
	sanitized = strings.ToUpper(sanitized)
	sanitized = nonAlphaNumericRegex.ReplaceAllString(sanitized, "")
	sanitized = strings.ReplaceAll(sanitized, " ", "")

	code := prefix + "_" + sanitized
	code = sanitizeToAllowedCharset(code)

	if len(code) > 30 {
		code = code[:30]
	}
	return code
}

// GenerateExpectedEpisodeCode sinh mã kỳ vọng cho một tập của series.
// Định dạng: {SERIES_CODE}_T{N}
func GenerateExpectedEpisodeCode(departmentCode, categoryCode, seriesTitle string, episodeNumber int) string {
	base := GenerateExpectedProgramCode(departmentCode, categoryCode, seriesTitle)
	return fmt.Sprintf("%s_T%d", base, episodeNumber)
}

// normalizeVietnamese loại bỏ dấu tiếng Việt.
// Xử lý Đ/đ tường minh (không phân rã bởi NFD), sau đó lược bỏ các dấu kết hợp.
func normalizeVietnamese(s string) string {
	s = strings.ReplaceAll(s, "Đ", "D")
	s = strings.ReplaceAll(s, "đ", "d")

	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	result, _, _ := transform.String(t, s)
	return result
}

// sanitizeToAllowedCharset chỉ giữ lại [A-Z0-9_]
func sanitizeToAllowedCharset(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// NormalizeTitleForComparison chuẩn hóa title để kiểm tra khả năng tái dùng.
// Loại bỏ dấu tiếng Việt, viết hoa, bỏ ký tự không phải chữ-số.
func NormalizeTitleForComparison(title string) string {
	s := normalizeVietnamese(title)
	s = strings.ToUpper(s)
	s = nonAlphaNumericRegex.ReplaceAllString(s, "")
	return strings.ReplaceAll(s, " ", "")
}
