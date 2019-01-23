// Prints the content found at a URL
package main

import (
	"fmt"
	"io"
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

		_, err = io.Copy(os.Stdout, resp.Body)

		if err != nil {
			fmt.Fprintf(os.Stderr, "reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

// successful curl
// go run main.go https://github.com

// bad curl
// go run main.go https://123
