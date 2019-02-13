package main

import (
	"fmt"
	"unicode/utf8"
)

func rotate(s []byte, i int) []byte {
	a := make([]byte, len(s))
	copy(a, s[i:])
	copy(a[len(s)-i:], s[:i])
	return a
}

func reverse(s []byte) {
	runeCount := 0
	for len(s) > 0 {
		_, size := utf8.DecodeRune(s)
		s = s[size:]
		runeCount++
	}

	for i := 0; i < runeCount/2; i++ {

	}
}

func main() {
	a := []byte("hello 世界")
	reverse(a)
	fmt.Println(a)
}
