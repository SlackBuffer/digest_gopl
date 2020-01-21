// A minimal "echo" server; prints headers, form data
package main

import (
	"fmt"
	"log"
	"net/http"
)

// behind the scene, the sever runs the handler for each incoming request in a separate goroutine
// so that it can serve multiple request simultaneously
func main() {
	// A request for /count invokes counter and all others invoke handler (browser would fetch /favicon.icon additionally)
	// A handler pattern that ends with a slash matches any URL that has the pattern as a prefix
	http.HandleFunc("/", handler) // handles incoming URLs that begins with "/"
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)

	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}

	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)

	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

// go run main.go
// curl http://localhost:8000/?q=query
