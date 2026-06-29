package util

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

// GenerateDocumentCode sinh một document code duy nhất (bất biến, tối đa 100 ký tự).
//
// Định dạng: {prefix}-{YYYYMMDD}-{6char-random}
// Ví dụ: PROPOSAL-20260502-7K2NQ4
//
// Bên gọi truyền prefix theo document_type (vd PROPOSAL, SCRIPT_MASTER, FORMAT_BIBLE).
// Hậu tố ngẫu nhiên dùng base32 (không padding, viết hoa, chữ-số — an toàn URL, không có
// ký tự dễ nhầm) để tránh trùng. 6 ký tự = 32^6 ≈ 1 tỷ entropy mỗi
// (prefix, ngày) — xác suất trùng không đáng kể ở quy mô kỳ vọng.
//
// Repository ExistsByCode + ràng buộc UNIQUE khi Create là lưới an toàn cho
// trường hợp trùng hiếm gặp (bên gọi retry Create khi lỗi trùng key).
func GenerateDocumentCode(prefix string) string {
	day := TodayInVN().Format("20060102")
	suffix := randomSuffix(6)
	return fmt.Sprintf("%s-%s-%s", strings.ToUpper(prefix), day, suffix)
}

// GeneratePackageCode sinh một content_package code duy nhất.
// Định dạng: PKG-{YYYYMMDD}-{6char-random}
func GeneratePackageCode() string {
	day := TodayInVN().Format("20060102")
	return fmt.Sprintf("PKG-%s-%s", day, randomSuffix(6))
}

// randomSuffix trả về chuỗi ngẫu nhiên mã hóa base32 độ dài n.
// Dùng crypto/rand cho bảo mật; bảng chữ base32 loại bỏ ký tự dễ nhầm
// (không có 0/O, 1/I/L) theo bảng chuẩn RFC 4648.
func randomSuffix(n int) string {
	// 5 byte mã hóa thành 8 ký tự base32; cần ceil(n*5/8) byte
	bytes := make([]byte, (n*5+7)/8)
	_, _ = rand.Read(bytes)
	enc := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)
	if len(enc) > n {
		enc = enc[:n]
	}
	return strings.ToUpper(enc)
}
