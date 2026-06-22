// Package util — strip 3-mark inline markdown helpers.
//
// Strips inline emphasis markers when computing derived fields (plain text,
// word count, search index) so emphasis chars don't pollute search matches
// or inflate word counts. Mirrors a front-end editor parser.
//
// Whitelist (must match the front-end parser exactly):
//   - **...** bold
//   - __...__ underline
//   - *...*  italic (single asterisk)
//   - \* \_ \\  backslash escape → literal char
//
// Single asterisks/underscores that aren't paired (e.g. "con_trỏ") are kept
// literal — matches FE behavior. Unclosed markers are kept literal too.
package util

import "strings"

// StripInlineMarkdown removes 3-mark inline markers (B/I/U) from input,
// returning the plain text suitable for derived fields. Backslash-escaped
// markers become literal chars.
//
// Algorithm mirrors FE parser: single-pass scanner with paired-marker check.
// O(n) time, O(n) output.
func StripInlineMarkdown(input string) string {
	if input == "" {
		return ""
	}
	return strings.Join(extractText(parseTokens(input)), "")
}

type markdownToken struct {
	isText   bool
	text     string
	mark     string // "bold" | "italic" | "underline"
	children []markdownToken
}

// markers ordered by length (longest first) so ** matches before *.
var markdownMarkers = []struct {
	token string
	mark  string
}{
	{"**", "bold"},
	{"__", "underline"},
	{"*", "italic"},
}

func parseTokens(input string) []markdownToken {
	out := []markdownToken{}
	var buf strings.Builder
	i := 0

	flush := func() {
		if buf.Len() > 0 {
			out = append(out, markdownToken{isText: true, text: buf.String()})
			buf.Reset()
		}
	}

	for i < len(input) {
		ch := input[i]

		// Backslash escape: \* / \_ / \\
		if ch == '\\' && i+1 < len(input) {
			next := input[i+1]
			if next == '*' || next == '_' || next == '\\' {
				buf.WriteByte(next)
				i += 2
				continue
			}
		}

		matched := false
		for _, m := range markdownMarkers {
			if !strings.HasPrefix(input[i:], m.token) {
				continue
			}
			closeIdx := findCloseMarker(input, i+len(m.token), m.token)
			if closeIdx == -1 {
				continue
			}
			if closeIdx == i+len(m.token) {
				// Empty mark → literal
				continue
			}
			flush()
			inner := input[i+len(m.token) : closeIdx]
			out = append(out, markdownToken{
				isText:   false,
				mark:     m.mark,
				children: parseTokens(inner),
			})
			i = closeIdx + len(m.token)
			matched = true
			break
		}
		if matched {
			continue
		}

		buf.WriteByte(ch)
		i++
	}
	flush()
	return out
}

func findCloseMarker(input string, start int, token string) int {
	i := start
	for i <= len(input)-len(token) {
		if input[i] == '\\' && i+1 < len(input) {
			i += 2
			continue
		}
		if strings.HasPrefix(input[i:], token) {
			return i
		}
		// Skip past longer asterisk/underscore runs when looking for shorter
		// marker (mirror FE behavior — token `*` skips `**`, etc.).
		if token == "*" && strings.HasPrefix(input[i:], "**") {
			i += 2
			continue
		}
		if token == "_" && strings.HasPrefix(input[i:], "__") {
			i += 2
			continue
		}
		i++
	}
	return -1
}

func extractText(tokens []markdownToken) []string {
	out := []string{}
	for _, t := range tokens {
		if t.isText {
			out = append(out, t.text)
			continue
		}
		out = append(out, extractText(t.children)...)
	}
	return out
}
