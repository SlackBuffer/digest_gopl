package main

import (
	"fmt"
	"math/bits"
)

func PopCount(x uint64) byte {
	var count byte
	for i, xx := 0, x; i < bits.Len64(x); i++ {
		count += byte(xx & 1)
		xx >>= 1
	}
	return count
}

func main() {
	fmt.Println(PopCount(2541234))
}
