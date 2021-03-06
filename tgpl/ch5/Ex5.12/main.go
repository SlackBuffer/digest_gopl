// prints the structure of the HTML tree in outline
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n) // pre is called before a node's children are visited
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n) // post is called after
	}
}

/* func outline(stack []string, n *html.Node) {
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
} */

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
	}
	// outline(nil, doc)

	var depth int
	var pre, post func(n *html.Node)

	pre = func(n *html.Node) {
		if n.Type == html.ElementNode {
			fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
			depth++
		}
	}
	post = func(n *html.Node) {
		if n.Type == html.ElementNode {
			depth--
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}
	forEachNode(doc, pre, post)
}

// curl http://www.gopl.io/ | go run main.go
