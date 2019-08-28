package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown. Press return to abort.")
	// time.Tick returns a channel on which it sends events periodically
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		// the select statement cause each iteration of the loop to wait up to 1 second for an abort, but not longer
		select {
		case <-tick:
			// Do nothing
		case <-abort:
			fmt.Println("Launch aborted!")
			// When this return happens, it stops receives events from `tick`, but the ticker goroutine is still there, trying in vain to send on a channel from which no goroutine is receiving - a goroutine leak
			return
		}
	}
	fmt.Println("Launched")
}
