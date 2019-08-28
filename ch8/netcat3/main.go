package main

import (
	"fmt"
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
		fmt.Println("done from sub goroutine")
	}()
	// 传过去要 echo 的内容
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
	fmt.Println("really done")
}

/* When the user closes the standard input stream (ctrl+d), `mustCopy` returns and the main goroutine calls `conn.Close()`, closing both halves of the network connection.
Closing the write half (main goroutine) of the connection causes the server to see an end-of-file condition.
Closing the read half causes the background goroutine’s call to `io.Copy` to return a "read from closed connection" error.

Before main returns, the background goroutine logs a message, then sends a value on the `done` channel. The main goroutine waits until it has received this value before returning.
As a result, the program always logs the "done" message before exiting.
*/

// ctrl+d sends EOF (Unix), to netcat3
