package util

import (
	"strings"
	"testing"
)

// FuzzStripInlineMarkdown kiểm chứng các invariant của StripInlineMarkdown
// với input ngẫu nhiên do Go sinh ra.
func FuzzStripInlineMarkdown(f *testing.F) {
	seeds := []string{
		"**bold**",
		"*i* __u__",
		`\* \_ \\`,
		"con_tro",
		"",
		"***",
		"a**b*c__d",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, in string) {
		out := StripInlineMarkdown(in)

		// Invariant 1: KHÔNG panic — chạy được tới đây nghĩa là đã qua.

		// Invariant 2: output không bao giờ dài hơn input (strip chỉ bỏ ký tự, không thêm).
		if len(out) > len(in) {
			t.Errorf("output dài hơn input: in=%q (%d) out=%q (%d)", in, len(in), out, len(out))
		}

		// Invariant 3: nếu output không còn ký tự mark (* _ \), strip lại phải ra y hệt (idempotent).
		if !strings.ContainsAny(out, `*_\`) {
			if again := StripInlineMarkdown(out); again != out {
				t.Errorf("không idempotent: %q -> %q -> %q", in, out, again)
			}
		}
	})
}
