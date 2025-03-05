//go:build !solution

package hogwarts

import "fmt"

const (
	white int8 = iota
	grey
	black
)

func topologySort(prereqs map[string][]string, vertex string, cash *map[string]int8, res *[]string) {
	(*cash)[vertex] = grey
	for _, v := range prereqs[vertex] {
		if (*cash)[v] == grey {
			panic("lol")
		}
		if (*cash)[v] != black {
			topologySort(prereqs, v, cash, res)
		}
	}
	(*cash)[vertex] = black
	*res = append(*res, vertex)
	return
}

func GetCourseList(prereqs map[string][]string) []string {
	cash := make(map[string]int8)
	res := make([]string, 0)
	for k := range prereqs {
		if cash[k] == white {
			topologySort(prereqs, k, &cash, &res)
		}
	}
	fmt.Println(res)
	return res
}
