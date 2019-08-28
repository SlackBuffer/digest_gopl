package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Commencing countdown.")

	// time.Tick returns a channel on which it sends events periodically
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick
	}
	fmt.Println("Launched")
}
