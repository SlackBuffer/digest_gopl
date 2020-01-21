package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func forEachNode(n *html.Node, nodeMap map[string][]*html.Node) {

	if n.Type == html.ElementNode {
		if _, ok := nodeMap[n.Data]; ok {
			nodeMap[n.Data] = append(nodeMap[n.Data], n)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, nodeMap)
	}

}

// modify later
func ElementByTagName(doc *html.Node, tag ...string) []*html.Node {
	if len(tag) == 0 {
		return []*html.Node{}
	}

	nodeMap := make(map[string][]*html.Node)
	for _, t := range tag {
		nodeMap[t] = []*html.Node{}
	}

	forEachNode(doc, nodeMap)
	fmt.Println(nodeMap)

	return []*html.Node{}
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
	}

	ElementByTagName(doc, "img")
	ElementByTagName(doc, "h1", "h2", "h3", "h4", "a")
}

// curl http://www.gopl.io/ | go run main.go
