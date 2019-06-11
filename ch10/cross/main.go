package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOOS, runtime.GOARCH)
}

// go build
// ./cross

// GOARCH=386 go build
// ./cross
