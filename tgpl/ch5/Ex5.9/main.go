// replaces each substring "$foo" within s by the text returned by f("foo")
package main

import (
	"fmt"
	"strings"
)

func expand(s string, f func(string) string) string {
	return strings.Replace(s, "foo", f("foo"), -1)
}

func foo(s string) string {
	var s1 []rune
	for _, r := range []rune(s) {
		s1 = append(s1, r+1)
	}
	return "_" + string(s1) + "_"
}

func main() {
	s := "abcfooacfoo"
	fmt.Println(expand(s, foo))
}
