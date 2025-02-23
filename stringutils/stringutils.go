package stringutils

import (
	"fmt"
	"strings"
	"unicode"
)

// escape invisible utf8 characters un the form of \uXXXX
func VisualizeInvisible(s string) string {
	var sb strings.Builder
	sb.Grow(int(float64(len(s)) * float64(1.5)))
	for _, r := range s {
		switch r {
		case ' ':
			sb.WriteByte(' ')
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		case '\v':
			sb.WriteString("\\v")
		case '\f':
			sb.WriteString("\\f")
		case '\b':
			sb.WriteString("\\b")
		default:
			if !unicode.IsPrint(r) {
				sb.WriteString(fmt.Sprintf("\\u%04X", r))
			} else {
				sb.WriteRune(r)
			}
		}
	}
	return sb.String()
}
