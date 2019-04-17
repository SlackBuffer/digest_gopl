package main

import "fmt"

func counter(out chan<- int) {
	fmt.Printf("%T\n", out)
	for x := 0; x < 100; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// implicit conversion occurs
	go counter(naturals)
	fmt.Printf("main: %T\n", naturals) // main: chan int
	go squarer(squares, naturals)
	printer(squares)
}
