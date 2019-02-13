package main

import "fmt"

func appendInt(x []int, y ...int) []int {
	var z []int
	zlen := len(x) + len(y)

	if zlen <= cap(x) {
		// there's room to grow; extends the slice
		z = x[:zlen]
	} else {
		// there's insufficient space; allocate a new array
		// grow by doubling avoids an excessive number of allocations and ensures that appending a single element takes constant time on average
		zcap := zlen
		// fmt.Println("no", zcap, 2*len(x))
		// first round -- len(x) == 0
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		// copies elements from one slice to another of the same type
		// `copy` returns the number of elements actually copied, which is the smaller of the two slice length
		copy(z, x)
	}

	// z[len(x)] = y
	copy(z[len(x):], y)
	return z
}

func main() {
	var x, y []int

	for i := 0; i < 10; i++ {
		y = appendInt(x, i)
		fmt.Printf("%d cap=%d\t%v\n", i, cap(y), y)
		x = y
	}
}
