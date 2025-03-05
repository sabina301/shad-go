//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	sb := strings.Builder{}
	for len(input) > 0 {
		r, size := utf8.DecodeLastRuneInString(input)
		sb.WriteRune(r)
		input = input[:len(input)-size]
	}
	return sb.String()
}
