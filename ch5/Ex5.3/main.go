package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func visit(n *html.Node) {
	if n.Type == html.TextNode {
		fmt.Println(n.Data)
	}

	if n.FirstChild != nil && n.FirstChild.Data != "script" {
		visit(n.FirstChild)
	}
	if n.NextSibling != nil && n.NextSibling.Data != "script" {
		visit(n.NextSibling)
	}
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	visit(doc)
}

// curl https://github.com | go run main.go >a.txt

// go build exercises-the_go_programming_language/ch1/Ex1.7.0
// ./Ex1.7.0 https://github.com | go run main.go
