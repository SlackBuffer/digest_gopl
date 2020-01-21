// uses a map whose keys represent the set of lines that have already appeared to ensure that subsequent occurrences not printed
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Go programmers often describe a map used in this fashion as a "set of strings"
	seen := make(map[string]bool)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if !seen[line] {
			seen[line] = true
			fmt.Println(line)
		}
	}

	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "dedup: %v\n", err)
		os.Exit(1)
	}
}
