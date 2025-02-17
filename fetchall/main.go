//go:build !solution

package main

import (
	"net/http"
	"os"
	"sync"
)

func doGet(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func main() {
	var wg sync.WaitGroup
	urls := os.Args[1:]
	for _, url := range urls {
		wg.Add(1)
		go doGet(url, &wg)
	}
	wg.Wait()
}
