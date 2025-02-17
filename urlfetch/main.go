//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]

	for _, url := range args {
		resp, err := http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		body, err := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}
