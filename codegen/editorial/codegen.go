package util

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

// GenerateDocumentCode produces an unique document code (immutable, max 100 chars).
//
// Format: {prefix}-{YYYYMMDD}-{6char-random}
// Example: PROPOSAL-20260502-7K2NQ4
//
// Caller passes prefix per document_type (e.g., PROPOSAL, SCRIPT_MASTER, FORMAT_BIBLE).
// Random suffix uses base32 (no padding, uppercase, alphanumeric — URL-safe, no
// confusable chars) for collision safety. 6 chars = 32^6 ≈ 1B entropy per
// (prefix, day) — collision probability negligible at expected scale.
//
// Repository ExistsByCode + Create UNIQUE constraint provides safety net for
// the rare collision case (caller retries Create on dup-key error).
func GenerateDocumentCode(prefix string) string {
	day := TodayInVN().Format("20060102")
	suffix := randomSuffix(6)
	return fmt.Sprintf("%s-%s-%s", strings.ToUpper(prefix), day, suffix)
}

// GeneratePackageCode produces an unique content_package code.
// Format: PKG-{YYYYMMDD}-{6char-random}
func GeneratePackageCode() string {
	day := TodayInVN().Format("20060102")
	return fmt.Sprintf("PKG-%s-%s", day, randomSuffix(6))
}

// randomSuffix returns base32-encoded random string of length n.
// Uses crypto/rand for security; base32 alphabet excludes confusable chars
// (no 0/O, 1/I/L) per RFC 4648 standard alphabet.
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
