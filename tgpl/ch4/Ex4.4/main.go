package main

import "fmt"

func rotate(s []int, i int) []int {
	// a := []int{}
	a := make([]int, len(s))
	copy(a, s[i:])
	copy(a[len(s)-i:], s[:i])
	return a
}

func main() {
	s := []int{0, 1, 2, 3, 4, 5, 6}
	fmt.Println(rotate(s, 3))
}
