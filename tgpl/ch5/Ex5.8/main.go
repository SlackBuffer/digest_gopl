// finds the first element with the specified id attribute
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func forEachNode(n *html.Node, id string) *html.Node {
	var node *html.Node // nil

	if n.Type == html.ElementNode {
		// fmt.Println(n.Data)
		for _, attr := range n.Attr {
			if attr.Key == "id" {
				if attr.Val == id {
					// fmt.Printf("%#v\n", n)
					return n
				} else {
					break
				}
			}
		}
	}

	for c := n.FirstChild; c != nil && node == nil; c = c.NextSibling {
		node = forEachNode(c, id)
	}

	return node
}

func ElementByID(doc *html.Node, id string) *html.Node {
	node := forEachNode(doc, id)
	return node
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
	}
	id := "toc"
	node := ElementByID(doc, id)
	if node == nil {
		fmt.Println("Such node doesn't exist.")
	} else {
		fmt.Printf("<%s id=%s>\n", node.Data, id)
	}
}

// curl http://www.gopl.io/ | go run main.go
