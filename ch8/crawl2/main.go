package main

import (
	"exercises-the_go_programming_language/ch5/links"
	"fmt"
	"log"
	"os"
)

// `tokens` is a counting semaphore used to enforces a limit of 20  concurrent requests (to avoid too parallel)
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token
	if err != nil {
		log.Print(err)
	}
	return list
}

// break out of the main loop when the worklist is empty and no crawl goroutines are active
func main() {
	worklist := make(chan []string)
	var n int // number of pending sends to `worklist`
	// start with the command-line argument
	n++
	go func() { worklist <- os.Args[1:] }()

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}
