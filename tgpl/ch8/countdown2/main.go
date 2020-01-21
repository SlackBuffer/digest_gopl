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

	select {
	// time.After immediately returns a channel, and starts a new goroutine that sends a single value on that channel after the specified time. This case is then satisfied (can read from that channel)
	case <-time.After(3 * time.Second):
		// Do nothing
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	fmt.Println("Launched")

	/* // `time.Tick` returns a channel on which it sends events periodically
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick
	} */
}
