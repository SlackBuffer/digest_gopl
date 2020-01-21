/*Package clargs measures the different in running time between the inefficient version and the one uses `strings.Join`*/
package clargs

import (
	"strings"
)

func printArgs1(args []string) {
	var s string
	sep := " "
	for _, arg := range args {
		s += arg + sep
	}
	// fmt.Println(s)
}

func printArgs2(args []string) {
	strings.Join(args, " ")
	// fmt.Println(strings.Join(args, " "))
}

/* func main() {
	args := os.Args[1:]

	printArgs2(args)
	// start := time.Now()
	// fmt.Printf("+=:\t\t %d nanoseconds elapsed\n", time.Since(start).Nanoseconds())

} */
