package main

import (
	"digest_gopl/ch5/links"
	"fmt"
	"log"
	"os"
)

func main() {
	worklist := make(chan []string)  // list of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	go func() { worklist <- os.Args[1:] }()

	// 先创建 20 个消费者待命
	// create 20 crawler goroutines to fetch each unseen link
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				
				// Links found by crawl are sent to the worklist from a dedicated goroutine to avoid deadlock
				// 若不起协程，一个 worklist 有超过 20 个 link，前 20 个占据 20 个消费者；第 21 个起， main 里 unseenLinks <- link 阻塞；20 个消费者协程 worklist <- foundLinks 写 worklist 通道后，main 里的下一轮 range worklist 由于前一轮还阻塞着无法执行，dead lock
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

// The crawler goroutines are fed by the same channel, unseenLinks.
// The main goroutine is responsible for de-douplicating items it receives from the worklist and then sending each unseen one over the unseenLinks channel to a crawl goroutine

// This program never terminates
