package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"unicode/utf8"
)

func comma(s string) string {
	var buf bytes.Buffer
	n := len(s)

	if n <= 3 {
		return s
	}

	buf.WriteString(s[:3])
	for i := 3; i < n; i++ {
		if s[i] == '.' {
			buf.WriteString(s[i:])
			break
		}

		if i%3 == 0 {
			buf.WriteByte(',')
		}
		r, _ := utf8.DecodeRuneInString(s[i:])
		// fmt.Printf("%c\n", r)
		buf.WriteRune(r)
	}
	return buf.String()
}

func main() {
	input := bufio.NewScanner(os.Stdin) // or file input
	for input.Scan() {
		fmt.Println(comma(input.Text()))
	}
}
