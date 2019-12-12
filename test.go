package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"time"
)

func main() {
	/* type A struct {
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
	// b, ok := w.(*bytes.Buffer) // failure: !ok, b == nil */

	fn := make(chan string)
	go func(ff chan<- string) {
		for _, f := range []string{"fn1", "fn2", "fn3"} {
			ff <- f
		}
		close(ff)
	}(fn)

	fmt.Println(pic(fn))

	// var as chan int
	// as := make(chan int)
	// go func() { as <- 2 }()
	// fmt.Println(<-as)

	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w)) // *os.File, not io.Writer

	v := reflect.ValueOf(3)

	fmt.Println(v)
	fmt.Println(v.String()) // <int Value>
	// <int Value>
	fmt.Println(reflect.ValueOf("qwr")) // qwr
	t := v.Type()
	fmt.Println(t) // int
}

func pic(fn <-chan string) string {
	str := make(chan string)
	var wg sync.WaitGroup
	for f := range fn {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			str <- f
		}(f)
	}

	go func() {
		wg.Wait()
		fmt.Println("zero")
		close(str)
	}()

	var ss string
	for s := range str {
		ss += s + "\n"
	}
	return ss
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
