package main

import "fmt"

// pc[i] is the population count of i
var pc [256]byte

/* precompute a table of results for each possible 8-bit value so that
PopCount function needn't take 64 steps but can just return the sum
of eight table lookups
*/
func init() {
	// i: index
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
		// fmt.Printf("%d\t%b\t%v\n", i, i, pc[i])

		/* fmt.Printf("%#v\n", i/2)	// 0 0 1 1 2 2 ... 127 127
		fmt.Printf("%v\n", byte(i&1)) // 0 1 0 1 0 1 ... */
	}
}

// PopCount return the population count (number of set bits) of x
// set bits: number of bits whose value is 1
func PopCount(x uint64) int {
	// truncate
	fmt.Printf("%b\t%b\t%b\n", x, byte(x), byte(x>>8))
	var res int
	for i := 0; i < 8; i++ {
		res += int(pc[byte(x>>(uint(i)*8))])
	}
	return res
}

func main() {
	fmt.Println(PopCount(2541234))
}
