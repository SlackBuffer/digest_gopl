package main

import (
	"exercises-the_go_programming_language/ch5/links"
	"fmt"
	"log"
	"os"
)

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)
	// must run in its own goroutine to avoid deadlock,
	// a stuck situation in which both the [ ] main goroutine and a crawler goroutine attempt to send to each other while neither is receiving
	go func() { worklist <- os.Args[1:] }()
	// 直接 block 掉了，不会往下跑
	// worklist <- os.Args[1:]

	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}
