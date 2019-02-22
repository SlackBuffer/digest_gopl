package main

import (
	"errors"
	"fmt"
)

func re() (s error) {
	defer func() {
		recover()
		s = errors.New("what")
	}()

	panic("returns value using panic instead of return")
}

func main() {
	fmt.Println(re())
}
