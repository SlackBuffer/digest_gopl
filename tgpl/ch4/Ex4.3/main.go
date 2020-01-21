package main

import "fmt"

const LENGTH = 7

func reverse(arr *[LENGTH]int) *[LENGTH]int {
	for i, j := 0, LENGTH-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func main() {
	a := [...]int{0, 1, 2, 3, 4, 5, 9}
	reverse(&a)
	fmt.Println(a)
}
