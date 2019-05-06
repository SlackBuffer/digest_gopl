package memo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	m := New(httpGetBody)
	urls := []string{"https://golang.org", "https://godoc.org", "https://play.golang.org", "http://gopl.io", "https://golang.org", "https://godoc.org", "https://play.golang.org", "http://gopl.io"}

	// wait util the last request is complete before returning
	var n sync.WaitGroup
	for _, url := range urls {
		n.Add(1)
		go func(url string) {
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))
			n.Done()
		}(url)
	}
	n.Wait()
	// unlikely to work correctly all the time, there exists data race

	/*
		// executes all calls to `Get` sequentially
		for _, url := range urls {
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))
		}
		// https://golang.org, 1.097030611s, 8158 bytes
		// https://godoc.org, 1.624867787s, 6805 bytes
		// https://play.golang.org, 1.354410201s, 6011 bytes
		// http://gopl.io, 2.233536468s, 4154 bytes
		// https://golang.org, 421ns, 8158 bytes
		// https://godoc.org, 176ns, 6805 bytes
		// https://play.golang.org, 168ns, 6011 bytes
		// http://gopl.io, 180ns, 4154 bytes
	*/
}

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// go run *.go
// go run -race *.go
