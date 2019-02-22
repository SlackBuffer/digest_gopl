package main

import "fmt"

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x)
	defer fmt.Printf("defer %d\n", x)
	f(x - 1)
}

func main() {
	f(3) // f(3) f(2) f(1) defer 1 defer 2 defer 3
}
