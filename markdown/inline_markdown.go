// Package util — các hàm lược bỏ 3 loại inline markdown.
//
// Lược bỏ các dấu nhấn mạnh inline khi tính các trường dẫn xuất (plain text,
// đếm từ, chỉ mục tìm kiếm) để ký tự nhấn mạnh không làm nhiễu kết quả tìm kiếm
// hay thổi phồng số từ. Mô phỏng đúng parser của editor phía front-end.
//
// Danh sách cho phép (phải khớp y hệt parser front-end):
//   - **...** đậm (bold)
//   - __...__ gạch chân (underline)
//   - *...*  nghiêng (italic, một dấu sao)
//   - \* \_ \\  escape bằng backslash → ký tự nguyên văn
//
// Dấu sao/gạch dưới đơn lẻ không thành cặp (vd "con_trỏ") được giữ nguyên văn
// — khớp hành vi FE. Dấu mở không có dấu đóng cũng giữ nguyên văn.
package util

import "strings"

// StripInlineMarkdown loại bỏ 3 loại dấu inline (đậm/nghiêng/gạch chân) khỏi input,
// trả về plain text dùng cho các trường dẫn xuất. Dấu được escape bằng backslash
// trở thành ký tự nguyên văn.
//
// Thuật toán mô phỏng parser FE: quét một lượt, kiểm tra dấu theo cặp.
// Thời gian O(n), output O(n).
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

// các dấu xếp theo độ dài (dài trước) để ** khớp trước *.
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

		// Escape backslash: \* / \_ / \\
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
				// Dấu rỗng → giữ nguyên văn
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
		// Bỏ qua các chuỗi sao/gạch dưới dài hơn khi đang tìm dấu ngắn hơn
		// (mô phỏng hành vi FE — dấu `*` bỏ qua `**`, v.v.).
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
