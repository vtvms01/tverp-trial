// Package readingpace ước lượng thời lượng đọc (nói) của text tiếng Việt.
//
// Chuẩn hóa text theo tiếng Việt + tính thời lượng dựa trên WPM (số từ mỗi phút).
//
// Công thức:
//
//	so_giay_noi  = so_tu_da_dieu_chinh / wpm_hieu_dung × 60
//	wpm_hieu_dung = wpm_cau_hinh × he_so_ngon_ngu × he_so_do_kho
//
// Chuẩn hóa tiếng Việt:
//   - "TP.HCM" → "Thanh pho Ho Chi Minh" (5 từ)
//   - "19h45" → "muoi chin gio bon muoi lam" (dạng đọc)
//   - "7,2%" → "bay phay hai phan tram" (dạng đọc)
package timing

import (
	"regexp"
	"strings"
)

// DefaultLanguageFactor — hệ số ngôn ngữ mặc định cho tiếng Việt.
const DefaultLanguageFactor = 1.0

// DefaultDifficultyFactor — hệ số độ khó mặc định cho text thường.
const DefaultDifficultyFactor = 1.0

// DifficultNameMultiplier — áp dụng khi từ khớp từ điển phát âm (bội số 1.2×).
const DifficultNameMultiplier = 1.2

// ─────────────────────────────────────────────────────────────────────────────
// Bảng mở rộng viết tắt tiếng Việt
// ─────────────────────────────────────────────────────────────────────────────

// vietnameseAbbreviations — các viết tắt thường gặp → dạng đọc.
// Key = mẫu viết tắt (sau chuẩn hóa); value = dạng đọc (chữ thường).
// Service áp dụng mở rộng TRƯỚC khi đếm từ.
var vietnameseAbbreviations = map[string]string{
	"TP.HCM": "Thanh pho Ho Chi Minh",
	"TP HCM": "Thanh pho Ho Chi Minh",
	"TPHCM":  "Thanh pho Ho Chi Minh",
	"TP.HN":  "Thanh pho Ha Noi",
	"VTV":    "ve te ve",
	"HTV":    "hat te ve",
	"VTC":    "ve te xe",
	"VN":     "Viet Nam",
	"GDP":    "ge de pe",
	"COVID":  "co vit",
	"USD":    "u es de",
	"VND":    "ve en de",
}

// digitInVN — tên đọc của chữ số đơn (0-9).
var digitInVN = map[byte]string{
	'0': "khong", '1': "mot", '2': "hai", '3': "ba", '4': "bon",
	'5': "nam", '6': "sau", '7': "bay", '8': "tam", '9': "chin",
}

// twoDigitInVN — các số hai chữ số đặc biệt 10-19.
var twoDigitInVN = map[string]string{
	"10": "muoi", "11": "muoi mot", "12": "muoi hai", "13": "muoi ba", "14": "muoi bon",
	"15": "muoi lam", "16": "muoi sau", "17": "muoi bay", "18": "muoi tam", "19": "muoi chin",
}

// ─────────────────────────────────────────────────────────────────────────────
// Các regex cho biến đổi sang dạng đọc
// ─────────────────────────────────────────────────────────────────────────────

var (
	// "19h45" → "muoi chin gio bon muoi lam"
	timeOfDayPattern = regexp.MustCompile(`\b(\d{1,2})h(\d{1,2})\b`)
	// "7,2%" → "bay phay hai phan tram" (khớp số,số%)
	percentDecimalPattern = regexp.MustCompile(`\b(\d+),(\d+)%`)
	// "7%" → "bay phan tram"
	percentPattern = regexp.MustCompile(`\b(\d+)%`)
	// Số nguyên đứng riêng (bắt các token số còn lại sau các biến đổi khác)
	standaloneIntPattern = regexp.MustCompile(`\b\d+\b`)
	// Thẻ HTML/markup <...>
	htmlTagPattern = regexp.MustCompile(`<[^>]+>`)
	// Cú pháp cue đôi [[...]] (phải xử lý TRƯỚC ngoặc đơn)
	cueDoubleBracketPattern = regexp.MustCompile(`\[\[[^\]]*\]\]`)
	// Cú pháp cue [...]
	cueBracketPattern = regexp.MustCompile(`\[[^\]]*\]`)
)

// ─────────────────────────────────────────────────────────────────────────────
// API công khai
// ─────────────────────────────────────────────────────────────────────────────

// ReadingPaceConfig chứa các hệ số riêng của người dẫn để tính WPM.
type ReadingPaceConfig struct {
	DefaultWPM        float64           // số từ mỗi phút cấu hình
	LanguageFactor    float64           // hệ số độ khó ngôn ngữ (tiếng Việt mặc định 1.0)
	DifficultyFactor  float64           // độ khó của text (mặc định 1.0)
	PronunciationDict map[string]string // presenter.pronunciation_dict — tên khó → dạng đọc
}

// PaceResult chứa kết quả trung gian + cuối cùng cho một khối text.
type PaceResult struct {
	NormalizedText    string  // text sau khi lược bỏ + mở rộng
	WordCount         int     // số từ thô
	AdjustedWordCount int     // tính cả bội số từ khó
	EffectiveWPM      float64 // wpm gốc × ngôn ngữ × độ khó
	SpokenSeconds     int     // floor(so_tu_da_dieu_chinh / wpm_hieu_dung * 60)
}

// CalculatePace tính thời lượng đọc (nói) cho một text.
//
// Các bước:
//  1. NormalizeText: lược bỏ HTML/cue/markup; mở rộng viết tắt + số + giờ + phần trăm
//  2. Đếm từ (tách theo khoảng trắng, sau chuẩn hóa)
//  3. Áp bội số từ khó (1.2×) cho từ có trong pronunciation_dict
//  4. Tính wpm_hieu_dung = gốc × ngôn ngữ × độ khó
//  5. so_giay_noi = floor(so_tu_da_dieu_chinh / wpm_hieu_dung × 60)
func CalculatePace(text string, cfg ReadingPaceConfig) PaceResult {
	if cfg.DefaultWPM <= 0 {

		cfg.DefaultWPM = 160.0 // nhịp đọc phát thanh tiếng Việt chuẩn
	}
	if cfg.LanguageFactor <= 0 {
		cfg.LanguageFactor = DefaultLanguageFactor
	}
	if cfg.DifficultyFactor <= 0 {
		cfg.DifficultyFactor = DefaultDifficultyFactor
	}

	normalized := NormalizeText(text)
	words := splitWords(normalized)
	wordCount := len(words)

	// Bội số từ khó — +0.2 cho mỗi từ khớp trong pronunciation_dict
	adjustedWC := float64(wordCount)
	if len(cfg.PronunciationDict) > 0 {
		extra := 0.0
		for _, w := range words {
			if _, ok := cfg.PronunciationDict[w]; ok {
				extra += (DifficultNameMultiplier - 1)
			}
		}
		adjustedWC += extra
	}

	effectiveWPM := cfg.DefaultWPM * cfg.LanguageFactor * cfg.DifficultyFactor
	if effectiveWPM <= 0 {
		effectiveWPM = cfg.DefaultWPM
	}

	spokenSecs := int((adjustedWC / effectiveWPM) * 60.0)
	if spokenSecs < 0 {
		spokenSecs = 0
	}

	return PaceResult{
		NormalizedText:    normalized,
		WordCount:         wordCount,
		AdjustedWordCount: int(adjustedWC + 0.5),
		EffectiveWPM:      effectiveWPM,
		SpokenSeconds:     spokenSecs,
	}
}

// NormalizeText lược bỏ markup + mở rộng viết tắt + số sang dạng đọc.
//
// Thứ tự thao tác quan trọng:
//  1. Lược bỏ thẻ HTML
//  2. Lược bỏ cú pháp cue [[...]] rồi [...]
//  3. Mở rộng viết tắt đã đăng ký (TP.HCM → Thanh pho Ho Chi Minh) TRƯỚC khi mở rộng số
//  4. Mở rộng mẫu giờ (19h45 → dạng đọc)
//  5. Mở rộng phần trăm thập phân (7,2% → dạng đọc)
//  6. Mở rộng phần trăm (7% → dạng đọc)
//  7. Mở rộng số nguyên đứng riêng
//  8. Gộp khoảng trắng thừa
func NormalizeText(text string) string {
	if text == "" {
		return ""
	}
	out := text

	// Lược bỏ markup
	out = htmlTagPattern.ReplaceAllString(out, " ")
	out = cueDoubleBracketPattern.ReplaceAllString(out, " ")
	out = cueBracketPattern.ReplaceAllString(out, " ")

	// Viết tắt tiếng Việt TRƯỚC khi mở rộng số (tránh "TP.HCM" bị tách giữa mẫu)
	for abbrev, spoken := range vietnameseAbbreviations {
		out = strings.ReplaceAll(out, abbrev, " "+spoken+" ")
	}

	// Giờ "19h45" → "muoi chin gio bon muoi lam"
	out = timeOfDayPattern.ReplaceAllStringFunc(out, func(match string) string {
		parts := timeOfDayPattern.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		hour := expandNumber(parts[1])
		minute := expandNumber(parts[2])
		return " " + hour + " gio " + minute + " "
	})

	// Phần trăm thập phân "7,2%" → "bay phay hai phan tram"
	out = percentDecimalPattern.ReplaceAllStringFunc(out, func(match string) string {
		trimmed := strings.TrimSuffix(match, "%")
		parts := strings.Split(trimmed, ",")
		if len(parts) != 2 {
			return match
		}
		integer := expandNumber(parts[0])
		decimal := expandNumber(parts[1])
		return " " + integer + " phay " + decimal + " phan tram "
	})

	// Phần trăm "7%" → "bay phan tram"
	out = percentPattern.ReplaceAllStringFunc(out, func(match string) string {
		trimmed := strings.TrimSuffix(match, "%")
		return " " + expandNumber(trimmed) + " phan tram "
	})

	// Số nguyên đứng riêng (bắt các token số còn lại)
	out = standaloneIntPattern.ReplaceAllStringFunc(out, expandNumber)

	// Gộp khoảng trắng
	out = strings.Join(strings.Fields(out), " ")
	return out
}

// expandNumber trả về dạng đọc của một chuỗi số (giới hạn 0-99).
// Số lớn hơn: dự phòng đọc từng chữ số (chấp nhận được ở Phase 1).
func expandNumber(s string) string {
	if v, ok := twoDigitInVN[s]; ok {
		return v
	}
	if len(s) == 1 {
		return digitInVN[s[0]]
	}
	// 20-99: ghép tên hàng chục + "muoi" + hàng đơn vị
	if len(s) == 2 {
		tens := s[0]
		units := s[1]
		var tensWord string
		switch tens {
		case '2':
			tensWord = "hai muoi"
		case '3':
			tensWord = "ba muoi"
		case '4':
			tensWord = "bon muoi"
		case '5':
			tensWord = "nam muoi"
		case '6':
			tensWord = "sau muoi"
		case '7':
			tensWord = "bay muoi"
		case '8':
			tensWord = "tam muoi"
		case '9':
			tensWord = "chin muoi"
		default:
			tensWord = digitInVN[tens]
		}
		if units == '0' {
			return tensWord
		}
		// Quy tắc tiếng Việt: "5" sau chữ số hàng chục → "lam" thay vì "nam"
		unitsWord := digitInVN[units]
		if units == '5' && tens != '1' {
			unitsWord = "lam"
		}
		return tensWord + " " + unitsWord
	}
	// Số từ 3 chữ số trở lên: dự phòng đọc từng chữ số (Phase 1 chấp nhận; Phase 2+ có thể mở rộng)
	var parts []string
	for i := 0; i < len(s); i++ {
		parts = append(parts, digitInVN[s[i]])
	}
	return strings.Join(parts, " ")
}

// splitWords tách text đã chuẩn hóa theo khoảng trắng; lọc bỏ chuỗi rỗng.
func splitWords(text string) []string {
	return strings.Fields(text)
}
