package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		echo(c, input.Text(), 1*time.Second)
	}
	// ignoring potential errors from `input.Err()`
	c.Close()
}

func main() {
	// creates an object that listens for incoming connections on a network port
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		// the listener's `Accept` method blocks until an incoming connection request request is made,
		// then returns a `net.Conn` representing the connection
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}

		// starts a new goroutine and keep the for loop going
		// otherwise only one connection can be handled,
		// because there's also a infinite loop at work
		// in `handleConn` function
		handleConn(conn) // handle one connection at a time
	}
}
