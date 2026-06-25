package util

import (
	"strings"
	"testing"
)

var benchSeeds = map[string]string{
	"plain":      "day la doan text khong dinh dang gi ca lap lai nhieu lan",
	"manyMarks":  "**dam** *nghieng* __gach chan__ **dam** *nghieng* __gach__ ",
	"manyEscape": `\*a\*b\*c\*d\*e\*f\*g\*h\*i\*j\*k\*l\*m\*n\*o\*p `,
}

func BenchmarkStripInlineMarkdown(b *testing.B) {
	for name, seed := range benchSeeds {
		text := strings.Repeat(seed, 200) // input dài
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = StripInlineMarkdown(text)
			}

		})
	}
}
