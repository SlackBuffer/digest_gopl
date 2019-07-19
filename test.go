package main

import (
	"fmt"
)

func main() {
	type A struct {
		name string
	}
	a := A{name: "ho"}
	b := struct{ name string }{name: "ho"}
	fmt.Println(a == b)
}
