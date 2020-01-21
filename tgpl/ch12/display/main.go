package main

import (
	"fmt"
	"os"
	"reflect"

	"digest_gopl/ch12/format"
)

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func main() {
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	Display("strangelove", strangelove)
	Display("os.Stderr", os.Stderr)

	// Apply Display to a reflect.Value and watch it traverse the internal representation of the type descriptor for *os.File
	Display("rV", reflect.ValueOf(os.Stderr))

	var i interface{} = 3
	// reflect.ValueOf always returns a Value of a concrete type since it extracts the contents of an interface value.
	Display("i", i)

	// 1. Display calls reflect.ValueOf(&i) returns a pointer to i, of kind Ptr.
	// 2. The switch case for Ptr calls Elem on this value, which returns a Value representing the variable i itself, of kind Interface.
	Display("&i", &i)
}

func Display(name string, x interface{}) {
	fmt.Printf("\nDislay %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))
}

func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, format.FormatAtom(key)), v.MapIndex(key))
		}

	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}

	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem())
		}

	// basic types, channels, funcs
	default:
		fmt.Printf("%s = %s\n", path, format.FormatAtom(v))
	}
}
