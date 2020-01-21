// Prints comamnd-line arguments
package main

/*
flag packageuses a program's command-line arguments to set the values of certain variables distributed throughout the program
*/
import (
	"flag"
	"fmt"
	"strings"
)

/*
flag.Bool creates a new flag variable of type bool
3 arguments: the name of the flag (`n`), the variable's default
value (`false`), and a message that will be printed if the user
provides an invalid argument, an invalid flag, or `-h` `-help`
*/

// -n causes the function to omit the trailing newline that would normally be printed
var n = flag.Bool("n", false, "omit trailing newline")

// -s causes the function to separate the output argument by the contents of the string `sep` instead of the default single space
var sep = flag.String("s", "*", "separator")

func main() {
	/*
		Must call `flag.Parse` before the flags are used,
		to update the flag variables from their default values
		If `flag.Parse` encouters an error, it prints a usage message
		and call os.Exit(2) to terminate the program
	*/
	flag.Parse()

	fmt.Print(strings.Join(flag.Args(), *sep))
	// The non-flag arguments are available from `flag.Args()` as a slice of strings

	if !*n {
		fmt.Println()
	}
}

// go run main.go a bc def
// go run main.go -s / a bc def
// go run main.go -n a bc def
// go run main.go -h
