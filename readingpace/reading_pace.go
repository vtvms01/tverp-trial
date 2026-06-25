// Package readingpace estimates spoken duration of Vietnamese text.
//
// Vietnamese-aware text normalization + WPM-based duration calculation.
//
// Formula:
//
//	spoken_text_seconds = adjusted_word_count / effective_wpm × 60
//	effective_wpm       = configured words-per-minute × language_factor × difficulty_factor
//
// Vietnamese normalization:
//   - "TP.HCM" → "Thanh pho Ho Chi Minh" (5 words)
//   - "19h45" → "muoi chin gio bon muoi lam" (spoken form)
//   - "7,2%" → "bay phay hai phan tram" (spoken form)
package timing

import (
	"regexp"
	"strings"
)

// DefaultLanguageFactor for Vietnamese.
const DefaultLanguageFactor = 1.0

// DefaultDifficultyFactor for normal text.
const DefaultDifficultyFactor = 1.0

// DifficultNameMultiplier — applied when word matches the pronunciation dictionary.
const DifficultNameMultiplier = 1.2

// ─────────────────────────────────────────────────────────────────────────────
// Vietnamese abbreviation expansion table
// ─────────────────────────────────────────────────────────────────────────────

// vietnameseAbbreviations — common spoken-form expansions.
// Key = abbreviation pattern (post-normalize); value = spoken form (lowercase).
// Service applies expansion BEFORE word counting.
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

// digitInVN — single-digit number names for spoken-form expansion (0-9).
var digitInVN = map[byte]string{
	'0': "khong", '1': "mot", '2': "hai", '3': "ba", '4': "bon",
	'5': "nam", '6': "sau", '7': "bay", '8': "tam", '9': "chin",
}

// twoDigitInVN — two-digit numbers 10-19 special cases.
var twoDigitInVN = map[string]string{
	"10": "muoi", "11": "muoi mot", "12": "muoi hai", "13": "muoi ba", "14": "muoi bon",
	"15": "muoi lam", "16": "muoi sau", "17": "muoi bay", "18": "muoi tam", "19": "muoi chin",
}

// ─────────────────────────────────────────────────────────────────────────────
// Regex patterns for spoken-form transformations
// ─────────────────────────────────────────────────────────────────────────────

var (
	// "19h45" → "muoi chin gio bon muoi lam"
	timeOfDayPattern = regexp.MustCompile(`\b(\d{1,2})h(\d{1,2})\b`)
	// "7,2%" → "bay phay hai phan tram" (matches digit,digit%)
	percentDecimalPattern = regexp.MustCompile(`\b(\d+),(\d+)%`)
	// "7%" → "bay phan tram"
	percentPattern = regexp.MustCompile(`\b(\d+)%`)
	// Standalone integer (catches remaining numeric tokens after other transforms)
	standaloneIntPattern = regexp.MustCompile(`\b\d+\b`)
	// HTML/markup tags <...>
	htmlTagPattern = regexp.MustCompile(`<[^>]+>`)
	// Double cue syntax [[...]] (must come BEFORE single bracket)
	cueDoubleBracketPattern = regexp.MustCompile(`\[\[[^\]]*\]\]`)
	// Cue syntax [...]
	cueBracketPattern = regexp.MustCompile(`\[[^\]]*\]`)
)

// ─────────────────────────────────────────────────────────────────────────────
// Public API
// ─────────────────────────────────────────────────────────────────────────────

// ReadingPaceConfig captures presenter-specific factors for WPM calculation.
type ReadingPaceConfig struct {
	DefaultWPM        float64           // configured words-per-minute
	LanguageFactor    float64           // language difficulty multiplier (Vietnamese default 1.0)
	DifficultyFactor  float64           // text difficulty (default 1.0)
	PronunciationDict map[string]string // presenter.pronunciation_dict — difficult names → spoken form
}

// PaceResult holds intermediate + final timing for a text block.
type PaceResult struct {
	NormalizedText    string  // post-strip post-expansion
	WordCount         int     // raw word count
	AdjustedWordCount int     // accounts for difficult-name multiplier
	EffectiveWPM      float64 // presenter base × language × difficulty
	SpokenSeconds     int     // floor(adjusted_word_count / effective_wpm * 60)
}

// CalculatePace computes spoken duration for a text.
//
// Steps:
//  1. NormalizeText: strip HTML/cue/markup; expand Vietnamese abbreviations + numbers + time-of-day + percent
//  2. Count words (whitespace split, after normalize)
//  3. Apply difficult-name multiplier (1.2×) for words in pronunciation_dict
//  4. Compute effective_wpm = base × language × difficulty
//  5. spoken_seconds = floor(adjusted_word_count / effective_wpm × 60)
func CalculatePace(text string, cfg ReadingPaceConfig) PaceResult {
	if cfg.DefaultWPM <= 0 {

		cfg.DefaultWPM = 160.0 // canonical Vietnamese broadcast pace
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

	// Difficult-name multiplier — +0.2 per word matched in pronunciation_dict
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

// NormalizeText strips markup + expands Vietnamese abbreviations + spoken-form numbers.
//
// Order of operations matters:
//  1. Strip HTML tags
//  2. Strip cue syntax [[...]] then [...]
//  3. Expand registered abbreviations (TP.HCM → Thanh pho Ho Chi Minh) BEFORE number expansion
//  4. Expand time-of-day patterns (19h45 → spoken)
//  5. Expand percent-decimal (7,2% → spoken)
//  6. Expand percent (7% → spoken)
//  7. Expand standalone integers
//  8. Collapse extra whitespace
func NormalizeText(text string) string {
	if text == "" {
		return ""
	}
	out := text

	// Strip markup
	out = htmlTagPattern.ReplaceAllString(out, " ")
	out = cueDoubleBracketPattern.ReplaceAllString(out, " ")
	out = cueBracketPattern.ReplaceAllString(out, " ")

	// Vietnamese abbreviations BEFORE number expansion (avoids "TP.HCM" being split mid-pattern)
	for abbrev, spoken := range vietnameseAbbreviations {
		out = strings.ReplaceAll(out, abbrev, " "+spoken+" ")
	}

	// Time-of-day "19h45" → "muoi chin gio bon muoi lam"
	out = timeOfDayPattern.ReplaceAllStringFunc(out, func(match string) string {
		parts := timeOfDayPattern.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		hour := expandNumber(parts[1])
		minute := expandNumber(parts[2])
		return " " + hour + " gio " + minute + " "
	})

	// Percent-decimal "7,2%" → "bay phay hai phan tram"
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

	// Percent "7%" → "bay phan tram"
	out = percentPattern.ReplaceAllStringFunc(out, func(match string) string {
		trimmed := strings.TrimSuffix(match, "%")
		return " " + expandNumber(trimmed) + " phan tram "
	})

	// Standalone integers (catch remaining numeric tokens)
	out = standaloneIntPattern.ReplaceAllStringFunc(out, expandNumber)

	// Collapse whitespace
	out = strings.Join(strings.Fields(out), " ")
	return out
}

// expandNumber returns spoken form of a numeric string (limited to 0-99).
// Larger numbers: simple digit-by-digit fallback (acceptable Phase 1).
func expandNumber(s string) string {
	if v, ok := twoDigitInVN[s]; ok {
		return v
	}
	if len(s) == 1 {
		return digitInVN[s[0]]
	}
	// 20-99: combine tens-digit name + "muoi" + units
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
		// Vietnamese rule: "5" after a tens digit → "lam" not "nam"
		unitsWord := digitInVN[units]
		if units == '5' && tens != '1' {
			unitsWord = "lam"
		}
		return tensWord + " " + unitsWord
	}
	// 3+ digit numbers: digit-by-digit fallback (Phase 1 acceptable; Phase 2+ may expand)
	var parts []string
	for i := 0; i < len(s); i++ {
		parts = append(parts, digitInVN[s[i]])
	}
	return strings.Join(parts, " ")
}

// splitWords splits normalized text on whitespace; filters empty strings.
func splitWords(text string) []string {
	return strings.Fields(text)
}
