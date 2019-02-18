// prints the links in an HTML document read from standard input
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	/* for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	} */
	if n.FirstChild != nil {
		links = visit(links, n.FirstChild)
	}
	if n.NextSibling != nil {
		links = visit(links, n.NextSibling)
	}

	return links
}

func main() {
	// `html.Parse` reads a sequence of bytes, parses them, and returns the root of the HTML document tree, an `html.Node`
	// HTML has several kinds of nodes - text, comments ...
	// deal with element node of the form <name key='value'>
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findslinks1: %v\n", err)
		os.Exit(1)
	}

	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

// curl https://github.com | go run main.go

// go build exercises-the_go_programming_language/ch1/Ex1.7.0
// ./Ex1.7.0 https://github.com | go run main.go

/* package html

import "io"

type Node struct {
	Type                    NodeType
	Data                    string
	Attr                    []Attribute
	FirstChild, NextSibling *Node
}

type NodeType int32

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

type Attribute struct {
	Key, Val string
}

func Parse(r io.Reader) (*Node, error) */
