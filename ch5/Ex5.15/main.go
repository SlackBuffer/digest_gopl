package main

import (
	"errors"
	"fmt"
	"os"
)

func max(n ...int) (int, error) {
	if len(n) == 0 {
		err := errors.New("at least one argument is required for max()")
		return 0, err
	} else {
		max := n[0]
		for _, m := range n {
			if m > max {
				max = m
			}
		}
		return max, nil
	}
}
func max1(a int, n ...int) int {
	max := a
	for _, m := range n {
		if m > max {
			max = m
		}
	}
	return max
}

func main() {
	fmt.Println(max1(1, 2, 3))

	n, err := max()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(n)
}
