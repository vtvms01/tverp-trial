// Package util — code generation helpers.
//
// Code generator: PREFIX-YYYYMMDD-XXXXXX (crypto-random suffix).
// Format: {PREFIX}-{YYYYMMDD}-{6-char base32 random}
// Security: crypto/rand + base32 RFC 4648 (no padding, uppercase, no confusable chars 0/O/1/I/L)
package util

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

// GenerateStoryCode returns a story code: ST-YYYYMMDD-XXXXXX.
//
// Format: ST + TodayInVN().Format("20060102") + randomSuffix(6).
// Date in VN timezone (UTC+7).
// 6-char base32 random → ~1B entropy/day → retry on UNIQUE violation.
//
// NEVER REUSED even after soft-delete — audit identity vĩnh viễn.
func GenerateStoryCode() string {
	day := TodayInVN().Format("20060102")
	return fmt.Sprintf("ST-%s-%s", day, randomSuffix(6))
}

// randomSuffix returns base32-encoded random string of length n.
//
// Uses crypto/rand for security; base32 RFC 4648 alphabet (uppercase, no padding) excludes
// confusable chars (no 0/O, 1/I/L). 5 bytes encode to 8 base32 chars; ceil(n*5/8) bytes needed.
//
// Random base32 suffix; never reused even after soft-delete.
func randomSuffix(n int) string {
	// 5 bytes encode to 8 base32 chars; ceil(n*5/8) bytes needed
	bytes := make([]byte, (n*5+7)/8)
	_, _ = rand.Read(bytes)
	enc := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)
	if len(enc) > n {
		enc = enc[:n]
	}
	return strings.ToUpper(enc)
}
