package main

import (
	"fmt"
	"unicode"
)

// squashes each run of adjacent Unicode spaces into a single ASCII space
func squash(c []byte) []byte {
	b := c[:0]
	for i, v := range c {
		if i == 0 {
			b = append(b, v)
			continue
		}
		if !unicode.IsSpace(rune(v)) {
			b = append(b, v)
		} else if !unicode.IsSpace(rune(b[len(b)-1])) {
			b = append(b, v)
		}
	}
	return b
}

func main() {
	s := "  ab  c   d d f"
	fmt.Println(string(squash([]byte(s))))
}
