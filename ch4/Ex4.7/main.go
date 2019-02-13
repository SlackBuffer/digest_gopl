package main

import (
	"fmt"
	"unicode/utf8"
)

func reverse(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func main() {
	s := []byte("你好 hello  世界")
	sz := 0
	for len(s[sz:]) > 0 {
		_, size := utf8.DecodeRune(s[sz:])
		if size > 1 {
			reverse(s[sz : sz+size])
		}
		sz += size
	}
	reverse(s)
	fmt.Println(string(s))
}
