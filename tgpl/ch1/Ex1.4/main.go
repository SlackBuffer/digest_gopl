// Prints the names of all files in which each duplicated line occurs
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]

	if len(files) == 0 {
		// handle inputs from standard input
		countLines(os.Stdin, counts, "")
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			// fmt.Printf("%T\n", arg) // filename, string
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
			countLines(f, counts, arg)
			f.Close()
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int, filename string) {
	input := bufio.NewScanner(f)
	if filename == "" { // input from os.Stdin, not file, thus no filename
		for input.Scan() {
			counts[input.Text()]++
		}
	} else {
		for input.Scan() {
			counts[input.Text()+" - "+filename]++
		}
	}
}

// go run main.go a.txt b.txt
