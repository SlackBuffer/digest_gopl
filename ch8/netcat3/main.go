package main

import (
	"io"
	"log"
	"net"
	"os"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

/* When the user closes the standard input stream, `mustCopy` returns and the main goroutine
calls `conn.Close()`, closing both halves of the network connection.
Closing the write half of the connection causes the server to see an end-of-file condition.
Closing the read half causes the background goroutine’s call to `io.Copy` to return a "read from closed connection" er ror */
