// A minimal "echo" server; prints incoming url path
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler) // handles incoming URLs that begins with "/"
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// run in two terminals
// go run main.go
// 2019/01/23 17:52:01 listen tcp 127.0.0.1:8000: bind: address already in use
// exit status 1
