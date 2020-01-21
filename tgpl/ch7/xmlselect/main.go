// Extract and print the text found beneath certain elements in an XML document tree.
// Do its job in a single pass over the input without ever materializing the tree.
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	dec := xml.NewDecoder(os.Stdin)

	var stack []string // stack of element names

	// The API guarantees that the sequence of StartElement and EndElement tokens will be properly matched, even in ill-formated documents.
	// Comments are ignored.
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok.Name.Local) // push
		case xml.EndElement:
			stack = stack[:len(stack)-1] // pop
		case xml.CharData: // e.g., <p>CharData</p>
			// prints the text only if the stack contains all the elements named by the command-line arguments, in order.
			if containsAll(stack, os.Args[1:]) {
				fmt.Printf("%s: %s\n", strings.Join(stack, " "), tok)
			}
		}
	}
}

// reports whether x contains the elements of y, in order.
func containsAll(x, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0] == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}

/*
go get gopl.io/ch1/fetch
go build gopl.io/ch1/fetch
./fetch http://www.w3.org/TR/2006/REC-xml11-20060816 |
./main div div h2
*/
