// a TCP server that periodically writes the time
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	// creates an object that listens for incoming connections on a network port
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		// the listener's `Accept` method blocks (waiting for a connection to come in) until an incoming connection request request is made,
		// then returns a `net.Conn` representing the connection
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}

		// The second must wait until the first client is finished because the server is sequential
		// Adding the `go` keyword causes each call to `handleConn` to run in its own goroutine
		// and let the main goroutine just handles the Accept thing

		// starts a new goroutine and keep the `for` loop going
		// otherwise only one connection can be handled,
		// because there's also an **infinite loop** at work
		// in `handleConn` function. The control won't go back
		// from `handleConn` to `Accept` until after one connection is done
		go handleConn(conn)
	}
}

// handles one complete client connection
func handleConn(c net.Conn) {
	defer c.Close()
	// the loop ends when the write fails, at which point
	// `handleConn` closes its side of the connection using
	// a deferred call to `Close` and goes back to waiting
	// for another request
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

// go run main.go
// "netcat"
// nc localhost 8000

// can also use telnet
