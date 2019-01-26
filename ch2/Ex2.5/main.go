package main

import (
	"fmt"
)

func PopCount(x uint64) byte {
	var count byte
	for {
		fmt.Printf("%8b\n", x)
		if x = x & (x - 1); x != 0 {
			fmt.Printf("%8b\n\n", x)
			count += 1
		} else {
			count += 1
			break
		}
	}
	return count
}

func main() {
	fmt.Println(PopCount(2541234))
}
