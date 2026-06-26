package util

import "testing"

func FuzzParseDateInVN(f *testing.F) {
	for _, s := range []string{"2026-01-21", "2024-02-29", "2026-13-99", "", "garbage", "0000-00-00"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, in string) {
		// Invariant 1: input rác KHÔNG được panic.
		result, err := ParseDateInVN(in)

		// Invariant 2: nếu parse thành công, round-trip phải parse lại được.
		if err == nil {
			formatted := result.Format("2006-01-02")
			if _, err2 := ParseDateInVN(formatted); err2 != nil {
				t.Errorf("round-trip fail: %q -> %q -> lỗi %v", in, formatted, err2)
			}
		}
	})
}
