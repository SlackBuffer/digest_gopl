package format

import (
	"reflect"
	"strconv"
)

/* func main() {
	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	fmt.Println(format.Any(x)) // "1"
	fmt.Println(format.Any(d)) // "1"

	fmt.Println(format.Any([]int64{x}))         // "[]int64 0x8202b87b0"
	fmt.Println(format.Any([]time.Duration{d})) // "[]time.Duration 0x8202b87e0"

	fmt.Println(format.Any(struct{ Name string }{"ho"})) // struct { Name string } value
} */

// Any formats any value as a string.
func Any(value interface{}) string {
	return FormatAtom(reflect.ValueOf(value))
}

/*
formatAtom formats a value without inspecting its internal structure.
	formatAtom treats each value as an indivisible thing with no internal structure.
 	Since Kind is concerned only with the underlying representation, format.Any works for named types too.
*/
func FormatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)

	// ...floating-point and complex cases omitted for brevity...

	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())

	// For reference types (channels, functions, pointers, slices, and maps), it prints the type and the reference address in hexadecimal.
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)

	// For aggregate types (structs and arrays) and interfaces it prints only the type of the value.
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}
