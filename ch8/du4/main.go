// The `du1` command computes the disk usage of the files in a directory
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var verbose = flag.Bool("v", false, "show verbose progress messages")

// a counting semaphore for limiting concurrency in dirents
var sema = make(chan struct{}, 200)

var done = make(chan struct{})

func main() {
	// cancel traversal when input is detect
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(done)
	}()

	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	var nfiles, nbytes int64
loop:
	for {
		select {
		case <-done:
			// drain filesizes to allow existing goroutines to finish (阻塞的数据不再写入)
			// drain the fileSizes channel, discarding all values until the channel is closed, before returns
			// does this to ensure any active calls to `walkDir` can run to completion without getting getting stuck sending to filesSizes
			for range fileSizes {
				// Do nothing
			}
			panic("cancelled")
			// return
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes has been closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}

	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// recursively walks the file tree rooted at `dir` and sends the size of each found file on fileSizes
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()

	// turns all goroutines created after cancellation into on-ops
	if cancelled() {
		return
	}

	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size() // in bytes
		}
	}
}

// returns the entries of directory dir
func dirents(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: // acquire token
	case <-done:
		return nil // cancelled
	}
	defer func() { <-sema }() // release token

	// `ReadDir` returns the same information that a call to `os.Stat` returns for a single file
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
	}
	return entries
}

// sudo go run main.go -v /Users/slackbuffer/Desktop
// sudo go run main.go -v $HOME /usr /bin /etc

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
