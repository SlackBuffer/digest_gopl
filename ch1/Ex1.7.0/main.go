// Prints the content found at a URL
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1) // cause the process to exit with a status code 1
		}

		// Body field is the server response as a readable stream
		b, err := ioutil.ReadAll(resp.Body)
		// Body stream is closed to avoid leaking resources
		resp.Body.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
}

// successful curl
// go run main.go https://github.com

// bad curl
// go run main.go https://123
