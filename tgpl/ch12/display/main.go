package main

import (
	"fmt"
	"os"
	"reflect"

	"digest_gopl/tgpl/ch12/format"
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
		// Len returns v's length.
		// It panics if v's Kind is not Array, Chan, Map, Slice, or String.
		for i := 0; i < v.Len(); i++ {
			// Index returns v's i'th element
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}

	case reflect.Struct:
		// NumField returns the number of fields in the struct v.
		for i := 0; i < v.NumField(); i++ {
			// Field returns a struct type's i'th field.
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}

	case reflect.Map:
		// MapKeys returns a slice containing all the keys present in the map,
		// in unspecified order.
		for _, key := range v.MapKeys() {
			// MapIndex returns the value associated with key in the map v.
			display(fmt.Sprintf("%s[%s]", path, format.FormatAtom(key)), v.MapIndex(key))
		}

	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			// Elem returns the value that the interface v contains
			// or that the pointer v points to.
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
/* 
Dislay strangelove (main.Movie):
strangelove.Title = "Dr. Strangelove"
strangelove.Subtitle = "How I Learned to Stop Worrying and Love the Bomb"
strangelove.Year = 1964
strangelove.Color = false
strangelove.Actor["Dr. Strangelove"] = "Peter Sellers"
strangelove.Actor["Grp. Capt. Lionel Mandrake"] = "Peter Sellers"
strangelove.Actor["Pres. Merkin Muffley"] = "Peter Sellers"
strangelove.Actor["Gen. Buck Turgidson"] = "George C. Scott"
strangelove.Actor["Brig. Gen. Jack D. Ripper"] = "Sterling Hayden"
strangelove.Actor["Maj. T.J. \"King\" Kong"] = "Slim Pickens"
strangelove.Oscars[0] = "Best Actor (Nomin.)"
strangelove.Oscars[1] = "Best Adapted Screenplay (Nomin.)"
strangelove.Oscars[2] = "Best Director (Nomin.)"
strangelove.Oscars[3] = "Best Picture (Nomin.)"
strangelove.Sequel = nil

Dislay os.Stderr (*os.File):
(*(*os.Stderr).file).pfd.fdmu.state = 0
(*(*os.Stderr).file).pfd.fdmu.rsema = 0
(*(*os.Stderr).file).pfd.fdmu.wsema = 0
(*(*os.Stderr).file).pfd.Sysfd = 2
(*(*os.Stderr).file).pfd.pd.runtimeCtx = 0
(*(*os.Stderr).file).pfd.iovecs = nil
(*(*os.Stderr).file).pfd.csema = 0
(*(*os.Stderr).file).pfd.isBlocking = 1
(*(*os.Stderr).file).pfd.IsStream = true
(*(*os.Stderr).file).pfd.ZeroReadIsEOF = true
(*(*os.Stderr).file).pfd.isFile = true
(*(*os.Stderr).file).name = "/dev/stderr"
(*(*os.Stderr).file).dirinfo = nil
(*(*os.Stderr).file).nonblock = false
(*(*os.Stderr).file).stdoutOrErr = true
(*(*os.Stderr).file).appendMode = false

Dislay rV (reflect.Value):
(*rV.typ).size = 8
(*rV.typ).ptrdata = 8
(*rV.typ).hash = 871609668
(*rV.typ).tflag = 1
(*rV.typ).align = 8
(*rV.typ).fieldAlign = 8
(*rV.typ).kind = 54
(*(*rV.typ).alg).hash = func(unsafe.Pointer, uintptr) uintptr 0x10532f0
(*(*rV.typ).alg).equal = func(unsafe.Pointer, unsafe.Pointer) bool 0x1002df0
(*(*rV.typ).gcdata) = 1
(*rV.typ).str = 8985
(*rV.typ).ptrToThis = 0
rV.ptr = unsafe.Pointer value
rV.flag = 22

Dislay i (int):
i = 3

Dislay &i (*interface {}):
(*&i).type = int
(*&i).value = 3 
*/