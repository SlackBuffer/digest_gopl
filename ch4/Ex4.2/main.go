// prints the SHA256 hash of its standard input by default
// supports a command-line flag to print the SHA384 or SHA512 hash instead
package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
)

var hash = flag.String("hash", "256", "hash function")

func main() {
	flag.Parse()
	// fmt.Println(*hash)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		switch *hash {
		case "384":
			fmt.Printf("%x\n", sha512.Sum384([]byte(input.Text())))
		case "512":
			fmt.Printf("%x\n", sha512.Sum512([]byte(input.Text())))
		default:
			fmt.Printf("%x\n", sha256.Sum256([]byte(input.Text())))
		}
	}
}
