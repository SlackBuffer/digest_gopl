package main

import (
	"exercises-the_go_programming_language/ch5/links"
	"fmt"
	"log"
	"os"
)

// calls f for each item in the worklist
// any items returned by f are added to the worklist
// f is called at most once for each item
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil

		for _, item := range items {
			if !seen[item] {
				seen[item] = true

				// f(item) returns all URLs of the current HTML page
				worklist = append(worklist, f(item)...) // f(item)... causes all the items in the list returned by f to be appended to the worklist
			}
		}
	}
}

// returns all URLs of the current HTML page
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	// crawl the web breadth-first
	// starting from the command-line arguments
	// the process ends when all reachable web pages have been crawled or the memory of the computer is exhausted
	breadthFirst(crawl, os.Args[1:])
}
