// FILE: pkg/util/program_code.go
// Package textcode generates and normalizes content codes.
//
// Vietnamese diacritic handling: Đ/đ do NOT decompose under NFD, so they are
// replaced explicitly BEFORE Unicode normalization.
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

// GenerateExpectedProgramCode generates expected code for a program.
// Format: {DEPT}_{CAT}_{TITLE}, max 30 chars.
// Deterministic: same title always yields the same code.
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

// GenerateExpectedEpisodeCode generates expected code for a series episode.
// Format: {SERIES_CODE}_T{N}
func GenerateExpectedEpisodeCode(departmentCode, categoryCode, seriesTitle string, episodeNumber int) string {
	base := GenerateExpectedProgramCode(departmentCode, categoryCode, seriesTitle)
	return fmt.Sprintf("%s_T%d", base, episodeNumber)
}

// normalizeVietnamese removes Vietnamese diacritics.
// Handles Đ/đ explicitly (not decomposed by NFD), then strips combining marks.
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

// sanitizeToAllowedCharset keeps only [A-Z0-9_]
func sanitizeToAllowedCharset(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// NormalizeTitleForComparison normalizes a title for reuse compatibility check.
// Strips Vietnamese diacritics, uppercases, removes non-alphanumeric.
func NormalizeTitleForComparison(title string) string {
	s := normalizeVietnamese(title)
	s = strings.ToUpper(s)
	s = nonAlphaNumericRegex.ReplaceAllString(s, "")
	return strings.ReplaceAll(s, " ", "")
}
