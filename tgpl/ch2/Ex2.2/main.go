package main

import (
	"bufio"
	"exercises-the_go_programming_language/ch2/tempconv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if args, l := os.Args[1:], len(os.Args[1:]); l != 0 {
		for _, arg := range args {
			convert(arg)
		}
	} else {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			convert(input.Text())
		}
	}
}

func convert(s string) {
	t, err := strconv.ParseFloat(s, 64)

	if err != nil {
		fmt.Fprintf(os.Stderr, "cf: %v\n", err)
	}

	f := tempconv.Fahrenheit(t)
	c := tempconv.Celsius(t)
	k := tempconv.Kelvin(t)

	fmt.Printf("%s = %s, %s = %s, %s = %s\n", f, tempconv.FToC(f), c, tempconv.CToF(c), k, tempconv.KToC(k))
}

// go run main.go 212 32
