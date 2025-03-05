//go:build !solution

package spacecollapse

import "strings"

func CollapseSpaces(input string) string {
	sb := strings.Builder{}
	count := 0
	for _, s := range input {
		if s != ' ' && s != '\t' && s != '\n' && s != '\r' {
			sb.WriteRune(s)
			count = 0
		} else if s == ' ' && count == 0 {
			sb.WriteRune(' ')
			count = 1
		} else if (s == '\t' || s == '\n' || s == '\r') && count == 0 {
			sb.WriteRune(' ')
			count = 1
		}
	}
	return sb.String()
}
