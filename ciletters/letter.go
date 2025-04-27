//go:build !solution

package ciletters

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed template.txt
var s string

func cut8Byte(s string) string {
	r := []rune(s)
	count := 8
	if len(r) < 8 {
		count = len(r)
	}
	return string(r[:count])
}

func cut10Str(s string) string {
	sS := strings.Split(s, "\n")
	sb := strings.Builder{}
	count := 10
	if len(sS) < 10 {
		count = len(sS)
	}
	for i := len(sS) - count; i < len(sS); i++ {
		_, _ = sb.WriteString(sS[i])
		if i != len(sS)-1 {
			_, _ = sb.WriteString("\n            ")
		} else {
			_, _ = sb.WriteString("\n\n        ")
		}
	}
	str := sb.String()
	return str
}

func MakeLetter(n *Notification) (string, error) {
	funcMap := template.FuncMap{
		"cut8Byte": cut8Byte,
		"cut10Str": cut10Str,
	}
	tmpl := template.Must(template.New("letter").Funcs(funcMap).Parse(s))
	b := bytes.Buffer{}
	err := tmpl.Execute(&b, *n)
	if err != nil {
		return "", err
	}
	bStr := b.String()
	return bStr[:len(bStr)-9], nil
}
