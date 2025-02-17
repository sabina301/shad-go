//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	mapa := make(map[string]int)
	files := os.Args[1:]

	for i := range files {
		file0, err := os.Open(files[i])
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		defer file0.Close()
		scanner := bufio.NewScanner(file0)
		for scanner.Scan() {
			line := scanner.Text()
			mapa[line]++
		}
	}

	for k, v := range mapa {
		if v > 1 {
			fmt.Printf("%d	%s\n", v, k)
		}
	}
}
