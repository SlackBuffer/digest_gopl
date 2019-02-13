package main

import "fmt"

// returns a slice holing only the non-empty strings
/* func nonempty2(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
} */

func nonempty(strings []string) []string {
	out := strings[:0] // zero-length slice of original
	for _, s := range strings {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func main() {
	data := []string{"a", "", "b"}
	fmt.Printf("%q\n", nonempty(data)) // ["a" "b"]
	fmt.Printf("%q\n", data)           // ["a" "b" "b"]
}
