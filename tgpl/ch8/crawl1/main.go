package main

import (
	"digest_gopl/ch5/links"
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

// breadthFirst
func main() {
	worklist := make(chan []string)

	// The initial send of the command-line arguments to the worklist must run in its own goroutine to avoid deadlock,
	// a stuck situation in which both the main goroutine and a crawler goroutine attempt to send to each other while neither is receiving
	go func() { worklist <- os.Args[1:] }()
	// worklist <- os.Args[1:]	// 直接 block 掉了，不会往下跑

	seen := make(map[string]bool)
	for list := range worklist {
		// list 从 make(chan []string) 类型的通道取值，值是 []string
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

// go run main.go http://gopl.io/

// The program created so many network connections at once that it exceeded the per-process limit on the number o open files,
// causing operations such as a DNS lookups and calls toe net.Dial to start failing.

// The program is too parallel. Unbounded parallelism is rarely a good idea since there's always a limiting factor in the system, such as the number of CUP cores for compute-bound workloads, the number of spindles and heads for local disk I/O operations, the bandwidth of the network for streaming downloads, or the serving capacity of a web service.

// The program never terminates
