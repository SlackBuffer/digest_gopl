package main

import (
	"exercises-the_go_programming_language/ch7/tempconv1"
	"flag"
	"fmt"
)

var temp = tempconv1.CelsiusFlag("temp", 20.0, "the temperature")

func main() {
	flag.Parse()

	// `fmt` calls the `String` method
	fmt.Println(*temp)
}

// go build main.go
// ./main help
// ./main
// ./main -temp -18C
