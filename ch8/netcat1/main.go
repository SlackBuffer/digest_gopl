// a read-only TCP client
package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

// Reads data from the connection and writes it to the standard output until an end-of-file condition or an error occurs。
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

// run 2 clients (netcat1) at the same time
// the second client must wait until the first is finished
// because the server (clock1) is sequential
// it deals with 1 client at a time
