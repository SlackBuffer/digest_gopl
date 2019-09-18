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

var sema = make(chan struct{}, 200)

// create a cancellation channel on which on values are ever sent, but whose closure indicates that it's time for the program to stop what it's doing
// done 关闭后 <-done 的 case 条件满足
var done = make(chan struct{})

func main() {
	// Cancel traversal when input is detect.
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
	// It might be profitable to poll the cancellation status again within walkDir's loop, to avoid creating goroutines after teh cancellation event.
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
		// Before this case returns, it must first drain the fileSizes channel, discarding all values until the channel is closed.
		// It does this to ensure any active calls to walkDir (goroutine) can run to completion without getting getting stuck sending to fileSizes.
		case <-done:
			// drain filesizes to allow existing goroutines to finish
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

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()

	// This turns all goroutines created after cancellation into no-ops
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

// A little [ ] profiling of this program revealed that the bottleneck was the acquisition of a semaphore token in dirents.
// The select below make this operation cancellable and reduces the typical cancellation latency of the program from hundreds of milliseconds to tens
func dirents(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: // acquire token
	case <-done:
		return nil // cancelled
	}
	defer func() { <-sema }() // release token

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
