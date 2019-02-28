package main

import (
	"flag"
	"fmt"
	"time"
)

// Duration **defines** a `time.Duration` flag with specified name, default value, and usage string
// The return value is the **address** of a `time.Duration` variable that stores the **value of the flag**
// The flag accepts a value acceptable to `time.ParseDuration`
var period = flag.Duration("period", 1*time.Second, "sleep period")

func main() {
	flag.Parse()
	fmt.Printf("Sleeping for %v...", *period)
	time.Sleep(*period)
	fmt.Println()
}

// go run main.go
// go run main.go -period 50ms
// go run main.go -period 2m30s
// go run main.go -period 1.5h
