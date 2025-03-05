//go:build !solution

package varfmt

import (
	"fmt"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	count := 0
	sk1 := 0
	sb := strings.Builder{}
	sb.Grow(len(format))
	var u uint64 = 0
	for _, v := range format {
		if v == '{' {
			sk1 = 1
		} else if sk1 >= 1 && v != '}' {
			sk1 = 2
			u += u*10 + uint64(v-'0')
		} else if v == '}' && sk1 >= 1 {
			if sk1 == 1 {
				fmt.Fprintf(&sb, "%v", args[count])
			} else {
				fmt.Fprintf(&sb, "%v", args[u])
			}
			sk1 = 0
			u = 0
			count++
		} else {
			sb.WriteRune(v)
		}
	}
	return sb.String()
}
