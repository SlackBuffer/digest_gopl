package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
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

	if err := run(); err != nil {
		fmt.Println(err)
	}
	bc := new(bytes.Buffer)
	fmt.Printf("%T", bc)

	// File, Buffer are structs, concrete type
	var w io.Writer = os.Stdout // os.Stdout is of type `*os.File`
	f, _ := w.(*os.File)        // success: ok, f == os.Stdout
	fmt.Println(f == os.Stdout)
	// type `*os.File` is not same as type `*bytes.Buffer`
	// b, ok := w.(*bytes.Buffer) // failure: !ok, b == nil
}

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", e.When, e.What)
}
func run() error {
	type a int

	return &MyError{
		time.Now(),
		"it didn't work",
	}

}
