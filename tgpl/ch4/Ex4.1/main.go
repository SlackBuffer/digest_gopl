// counts the number of bits that are different in 2 SHA256 hases
package main

import (
	"crypto/sha256"
	"fmt"
)

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

		// fmt.Printf("%#v\n", i/2)      // 0 0 1 1 2 2 ... 127 127
		// fmt.Printf("%v\n", byte(i&1)) // 0 1 0 1 0 1 ...
	}
}

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))

	count := byte(0)
	for i, n := range c1 {
		t := n ^ c2[i]
		count += pc[t]
		// fmt.Printf("%x\t%x\n", n, c2[i])
	}

	fmt.Println(count)

	/* s1 := c1[:]
	s2 := c2[:]
	dst1 := make([]byte, hex.EncodedLen(len(s1)))
	dst2 := make([]byte, hex.EncodedLen(len(s2)))
	hex.Encode(dst1, s1)
	hex.Encode(dst2, s2) */

	/* r1 := fmt.Sprintf("%x", c1)
	r2 := fmt.Sprintf("%x", c2)
	fmt.Printf("%[1]s\n%[1]T\n", r1)
	fmt.Printf("%[1]s\n%[1]T\n", r2) */

}
