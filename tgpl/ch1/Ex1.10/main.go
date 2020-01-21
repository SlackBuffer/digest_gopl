// Fetches URLs in parallel and reports elapsed time and size; investgates caching
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)

	for _, url := range os.Args[1:] {
		go fetch(url, ch) // start a goroutine
	}

	// main does all the printing, ensures output from each goroutine is processed as a unit
	// without no danger of 2 goroutines finish at the same time (call print at the same time)
	for i, _ := range os.Args[1:] {
		res := strings.Split(<-ch, "----")
		err := ioutil.WriteFile("res"+strconv.Itoa(i)+".html", []byte(res[1]), 0644)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		fmt.Println(res[0])
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)

	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	nbytes := len(b)
	resp.Body.Close()

	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s----%s", secs, nbytes, url, string(b[:]))
}

// go run main.go https://github.com https://github.com
