package main

import (
	"bufio"
	"fmt"
	"strings"
)

type Counter struct {
	word, line int
}

func (c *Counter) Write(p []byte) (int, error) {
	input := bufio.NewScanner(strings.NewReader(string(p)))
	for input.Scan() {
		// c.line++
		(*c).line++

		s := bufio.NewScanner(strings.NewReader(input.Text()))
		s.Split(bufio.ScanWords)
		for s.Scan() {
			c.word++
		}
	}
	return 0, nil
}

func main() {
	s := "Now is the winter of our discontent,\nMade glorious summer by this sun of York.\n"

	var c Counter
	fmt.Fprintf(&c, "%s", s)
	fmt.Printf("%#v\n", c)
}
