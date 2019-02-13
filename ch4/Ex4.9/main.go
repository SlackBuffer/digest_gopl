// reports the frequency of each word in an input text file
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var counts = make(map[string]int)

func wordfreq(s string) {
	counts[s]++
}

func main() {
	file, err := os.Open("text.md")
	if err != nil {
		log.Fatal(err)
	}

	input := bufio.NewScanner(file)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		wordfreq(input.Text())
	}
	for k, v := range counts {
		fmt.Printf("%q: %d\n", k, v)
	}
}
