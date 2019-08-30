package memo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// It's unlikely to work correctly all the time (there exists data race).
// We may notice unexpected cache misses, or cache hits that return incorrect values, or even crashes.
// Worse, it's likely to work correctly some of the time, so we may not even notice that it has a problem (if we don't run it with -race flag).
func main() {
	m := New(httpGetBody)
	// URLs contain duplicates
	urls := []string{"https://golang.org", "https://godoc.org", "https://play.golang.org", "http://gopl.io", "https://golang.org", "https://godoc.org", "https://play.golang.org", "http://gopl.io"}

	// Wait util the last request is complete before returning
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

// go test -run=TestConcurrent -race -v digest_gopl/ch9/memo1

/*
WARNING: DATA RACE
Write at 0x00c0000aec00 by goroutine 32:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.12.7/libexec/src/runtime/map_faststr.go:202 +0x0
  digest_gopl/ch9/memo1.(*Memo).Get()
      /Users/slackbuffer/go/src/digest_gopl/ch9/memo1/memo.go:30 +0x1ce
  digest_gopl/ch9/memotest.Concurrent.func1()
      /Users/slackbuffer/go/src/digest_gopl/ch9/memotest/memotest.go:105 +0xc5

Previous write at 0x00c0000aec00 by goroutine 15:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.12.7/libexec/src/runtime/map_faststr.go:202 +0x0
  digest_gopl/ch9/memo1.(*Memo).Get()
      /Users/slackbuffer/go/src/digest_gopl/ch9/memo1/memo.go:30 +0x1ce
  digest_gopl/ch9/memotest.Concurrent.func1()
      /Users/slackbuffer/go/src/digest_gopl/ch9/memotest/memotest.go:105 +0xc5

Goroutine 32 (running) created at:
  digest_gopl/ch9/memotest.Concurrent()
      /Users/slackbuffer/go/src/digest_gopl/ch9/memotest/memotest.go:102 +0x10c
  digest_gopl/ch9/memo1_test.TestConcurrent()
      /Users/slackbuffer/go/src/digest_gopl/ch9/memo1/memo_test.go:23 +0xda
  testing.tRunner()
      /usr/local/Cellar/go/1.12.7/libexec/src/testing/testing.go:865 +0x163

Goroutine 15 (finished) created at:
  digest_gopl/ch9/memotest.Concurrent()
      /Users/slackbuffer/go/src/digest_gopl/ch9/memotest/memotest.go:102 +0x10c
  digest_gopl/ch9/memo1_test.TestConcurrent()
      /Users/slackbuffer/go/src/digest_gopl/ch9/memo1/memo_test.go:23 +0xda
  testing.tRunner()
	  /usr/local/Cellar/go/1.12.7/libexec/src/testing/testing.go:865 +0x163
*/
