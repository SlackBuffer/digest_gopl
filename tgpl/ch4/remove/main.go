package main

import "fmt"

// preserving the order
func remove(slice []int, i int) []int {
	fmt.Println(slice[i:], slice[i+1:])
	copy(slice[i:], slice[i+1:])
	fmt.Println(slice, slice[:len(slice)-1])
	fmt.Println()

	return slice[:len(slice)-1]
}

// not preserving the order
func remove2(slice []int, i int) []int {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func main() {
	s := []int{5, 6, 7, 8, 9}
	fmt.Println(remove(s, 2))
}
