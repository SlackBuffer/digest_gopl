package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func visit(eleType map[string]int, n *html.Node) map[string]int {
	if n.Type == html.ElementNode {
		eleType[n.Data]++
	}

	/* for c := n.FirstChild; c != nil; c = c.NextSibling {
		visit(eleType, c)
	} */

	if n.FirstChild != nil {
		visit(eleType, n.FirstChild)
	}
	if n.NextSibling != nil {
		visit(eleType, n.NextSibling)
	}

	return eleType
}

func main() {
	eleType := make(map[string]int)

	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	for key, value := range visit(eleType, doc) {
		fmt.Printf("%s\t%d\n", key, value)
	}
}

// curl https://github.com | go run main.go

// go build exercises-the_go_programming_language/ch1/Ex1.7.0
// ./Ex1.7.0 https://github.com | go run main.go
