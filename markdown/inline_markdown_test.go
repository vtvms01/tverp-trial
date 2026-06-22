package util

import "testing"

func TestStripInlineMarkdown(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"plain text", "Hello world", "Hello world"},
		{"bold", "**foo**", "foo"},
		{"italic", "*foo*", "foo"},
		{"underline", "__foo__", "foo"},
		{"mixed", "a **bold** b *it* c __u__ d", "a bold b it c u d"},
		{"nested cross-type", "**bold __under__**", "bold under"},
		{"escape asterisk", `\*\*1500\*\*`, "**1500**"},
		{"escape backslash", `a\\b`, `a\b`},
		{"unclosed bold", "**foo", "**foo"},
		{"empty mark literal", "****", "****"},
		{"vietnamese single underscore literal", "con_trỏ và biến_số", "con_trỏ và biến_số"},
		{"vietnamese inside bold", "**Mặt trời bé con**", "Mặt trời bé con"},
		{"XSS-style HTML chars passthrough literal", "<script>alert(1)</script>", "<script>alert(1)</script>"},
		{"multiline preserved", "line1\n**line2**\nline3", "line1\nline2\nline3"},
		{"multiple paragraphs with marks", "Hôm nay **TP.HCM** ghi nhận __3 ca__ mới", "Hôm nay TP.HCM ghi nhận 3 ca mới"},
		{"trailing escape backslash alone (not at end)", `foo\bar`, `foo\bar`},
		{"adjacent marks", "**a***b*", "ab"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := StripInlineMarkdown(tc.in)
			if got != tc.want {
				t.Errorf("StripInlineMarkdown(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
