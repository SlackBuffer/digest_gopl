// Prints command-line arguments
package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func main() {
	var s string
	sep := " "
	args := os.Args
	aSlice := make([]string, 1)

	fmt.Println("Type of `os.Args` is a slice: ", reflect.TypeOf(args) == reflect.TypeOf(aSlice))
	fmt.Println("Command name: ", args[0])

	// 1. a quadratic process that could be costly if the number of arguments is large
	for i := 1; i < len(args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println("1. ", s)

	// 2. same costly issue
	s = ""
	for _, arg := range args[1:] {
		s += sep + arg
	}
	fmt.Println("2. ", s)

	// 3. efficient
	fmt.Println("3. ", strings.Join(os.Args[1:], " "))

	// 4. don't care about format, just to see the values
	fmt.Println("4. ", os.Args[1:])
}

// go run main.go arg1 arg2
// go build && ./Ex1.1 arg1 arg2
