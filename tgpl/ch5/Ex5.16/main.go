package main

import (
	"fmt"
	"strings"
)

func join(sep string, str ...string) {
	fmt.Println(strings.Join(str, sep))
}

func main() {
	join("-", "a", "b", "c")
	join("-")
}
