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

	var ab struct{} // type struct{}, empty struct
	ab = struct{}{} // empty struct literal
	fmt.Printf("%#v\n", ab)
}
