// prints the structure of the HTML tree in outline
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

var depth int

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

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		s := n.Data
		for _, attr := range n.Attr {
			s += " " + attr.Key + "=" + "'" + attr.Val + "'"
		}
		if n.FirstChild != nil {
			fmt.Printf("%*s<%s>\n", depth*2, "", s)
		} else {
			fmt.Printf("%*s<%s/>\n", depth*2, "", s)
		}
		depth++
	} else if n.Type == html.TextNode {
		fmt.Printf("%*s%s\n", depth*2, "", n.Data)
		depth++
	} else if n.Type == html.CommentNode {
		fmt.Printf("%*s<!-- %s -->\n", depth*2, "", n.Data)
		depth++
	}
}
func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		if n.FirstChild != nil {
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	} else if n.Type == html.TextNode || n.Type == html.CommentNode {
		depth--
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
	forEachNode(doc, startElement, endElement)
}

// curl http://www.gopl.io/ | go run main.go
