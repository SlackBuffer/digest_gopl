// prints the links in an HTML document read from standard input
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func visit(links map[string]int, n *html.Node) map[string]int {
	if n.Type == html.ElementNode {
		// script, img: src; link, a: href
		if n.Data == "a" || n.Data == "link" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links[n.Data]++
					break
				}
			}
		} else if n.Data == "script" || n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					links[n.Data]++
					break
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visit(links, c)
	}

	return links
}

func main() {
	links := make(map[string]int)

	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findslinks1: %v\n", err)
		os.Exit(1)
	}

	for key, value := range visit(links, doc) {
		fmt.Printf("%s\t%d\n", key, value)
	}
}
