// Package util — các helper sinh mã.
//
// Bộ sinh mã: PREFIX-YYYYMMDD-XXXXXX (hậu tố ngẫu nhiên crypto).
// Định dạng: {PREFIX}-{YYYYMMDD}-{6 ký tự base32 ngẫu nhiên}
// Bảo mật: crypto/rand + base32 RFC 4648 (không padding, viết hoa, không ký tự dễ nhầm 0/O/1/I/L)
package util

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

// GenerateStoryCode trả về một story code: ST-YYYYMMDD-XXXXXX.
//
// Định dạng: ST + TodayInVN().Format("20060102") + randomSuffix(6).
// Ngày theo múi giờ VN (UTC+7).
// 6 ký tự base32 ngẫu nhiên → ~1 tỷ entropy/ngày → retry khi vi phạm UNIQUE.
//
// KHÔNG BAO GIỜ TÁI DÙNG kể cả sau soft-delete — định danh audit vĩnh viễn.
func GenerateStoryCode() string {
	day := TodayInVN().Format("20060102")
	return fmt.Sprintf("ST-%s-%s", day, randomSuffix(6))
}

// randomSuffix trả về chuỗi ngẫu nhiên mã hóa base32 độ dài n.
//
// Dùng crypto/rand cho bảo mật; bảng chữ base32 RFC 4648 (viết hoa, không padding) loại bỏ
// ký tự dễ nhầm (không có 0/O, 1/I/L). 5 byte mã hóa thành 8 ký tự base32; cần ceil(n*5/8) byte.
//
// Hậu tố base32 ngẫu nhiên; không bao giờ tái dùng kể cả sau soft-delete.
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
