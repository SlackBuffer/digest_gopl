// does an HTTP GET request for the HTML document url and
// returns the number of words and images in it
package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func countWords(s string) (words int) {
	input := bufio.NewScanner(strings.NewReader(s))
	input.Split(bufio.ScanWords)
	for input.Scan() {
		words++
	}
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	switch n.Type {
	case html.ElementNode:
		if n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					// fmt.Println(a.Val)
				}
			}
			images++
		}
	case html.TextNode:
		fmt.Println(n.Data)
		words += countWords(n.Data)
	}

	if n.FirstChild != nil {
		w, i := countWordsAndImages(n.FirstChild)
		words, images = words+w, images+i
	}
	if n.NextSibling != nil {
		w, i := countWordsAndImages(n.NextSibling)
		words, images = words+w, images+i
	}

	return
}

func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return // return 0 0 err;
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return // return 0 0 err;
	}

	words, images = countWordsAndImages(doc)

	return // return words, images, nil
}

func main() {
	for _, url := range os.Args[1:] {
		fmt.Println(CountWordsAndImages(url))
	}
}
