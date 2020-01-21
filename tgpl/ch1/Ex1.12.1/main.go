// A minimal "echo" server; prints incoming url path; count request counts
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

// behind the scene, the sever runs the handler for each incoming request in a separate goroutine
// so that it can serve multiple request simultaneously
func main() {
	// A request for /count invokes counter and all others invoke handler (browser would fetch /favicon.icon additionally)
	// A handler pattern that ends with a slash matches any URL that has the pattern as a prefix
	http.HandleFunc("/count", counter)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", handler) // handles incoming URLs that begins with "/"
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

//
