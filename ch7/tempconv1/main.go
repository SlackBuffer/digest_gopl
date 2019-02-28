package tempconv1

import (
	"exercises-the_go_programming_language/ch2/tempconv"
	"flag"
	"fmt"
)

// get `String` from tempconv.Celsius for free
type celsiusFlag struct {
	tempconv.Celsius
}

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)

	switch unit {
	case "C", "°C":
		f.Celsius = tempconv.Celsius(value)
		return nil

	case "F", "°F":
		f.Celsius = tempconv.FToC(tempconv.Fahrenheit(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

// defines a Celsius flag with the specified name, default value,
// and usage, and returns the address of the flag variable
// the flag argument must have a quantity and a unit, e.g., 100C
func CelsiusFlag(name string, value tempconv.Celsius, usage string) *tempconv.Celsius {
	// init the default value
	f := celsiusFlag{value}

	// `Var` defines a flag with the specified name and usage string
	// The type and value of the flag are represented by the first argument, of type `Value`, which typically holds a user-defined implementation of `Value`
	// For instance, the caller could create a flag that turns a comma-separated string into a slice of strings by giving the slice the methods of `Value`; in particular, `Set` would decompose the comma-separated string into the slice

	// adds the flag to the application's set of command-line flags - the global variable `flag.CommandLine`
	// assign a `*celsiusFlag` argument to a `flag.Value` parameter, causing the compiler to check that `*celsiusFlag` has the necessary methods
	flag.CommandLine.Var(&f, name, usage)

	fmt.Printf("%#v\n", f)

	// returns a pointer to the Celsius field embedded within the celsiusFlag variable f
	// Celsius field is the variable that will be updated within the `Set` method during flags process
	return &f.Celsius
}

// 命令行按格式传参后会调用对应类型的 flag 的 Set 函数
