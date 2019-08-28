package main

import "fmt"

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; ; x++ {
			naturals <- x
		}
	}()

	// Squarer
	go func() {
		for {
			x, ok := <-naturals
			if !ok {
				break
			}
			squares <- x * x
		}
		close(squares)
	}()

	for {
		fmt.Println(<-squares)
	}
}

// We've intentionally chosen very simple functions, though of course they are
// too computationally trivial to warrent their own goroutines in a realistic program.
// Pipelines like this may be found in long-running server programs where channels
// are used for lifelong communication between goroutines containing infinite loops.
