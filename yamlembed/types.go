package yamlembed

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

type Foo struct {
	A string
	p int64
}

type Bar struct {
	I      int64    `yaml:"i,omitempty"`
	B      string   `yaml:"b"`
	UpperB string   `yaml:"-"`
	OI     []string `yaml:"oi,omitempty"`
	F      []any    `yaml:"f"`
}

func (b *Bar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tB := struct {
		B  string   `yaml:"b"`
		OI []string `yaml:"oi"`
		F  []any    `yaml:"f"`
	}{}

	if err := unmarshal(&tB); err != nil {
		return err
	}
	b.B = tB.B
	b.UpperB = strings.ToUpper(tB.B)
	b.OI = tB.OI

	b.F = make([]any, len(tB.F))
	for i, v := range tB.F {
		b.F[i] = v
	}

	return nil
}
func (b *Bar) MarshalYAML() (interface{}, error) {
	type alias Bar
	tB := struct {
		B string `yaml:"b"`
		F string `yaml:"f"`
	}{
		B: b.B,
		F: fmt.Sprintf("[%v]", strings.Join(stringifySlice(b.F), ", ")),
	}

	out, err := yaml.Marshal(tB)
	if err != nil {
		return nil, err
	}
	return out, err
}

func stringifySlice(slice []any) []string {
	var result []string
	for _, v := range slice {
		result = append(result, fmt.Sprintf("%v", v))
	}
	return result
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}
