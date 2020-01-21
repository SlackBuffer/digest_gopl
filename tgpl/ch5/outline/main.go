// prints the structure of the HTML tree in outline
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data) // push tag
		fmt.Println(stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// When outline calls itself recursively, the callee receives a **copy** of stack

		// Although the callee may append elements to this slice, modifying its underlying
		// array and even allocating a new array, it doesn't modify the initial elements that
		// are visible to the caller, so when the function returns, the caller's stack is as
		// it was before the call
		outline(stack, c)
	}
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
	}
	outline(nil, doc)
}

// curl https://github.com | go run main.go
