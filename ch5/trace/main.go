package main

import (
	"log"
	"time"
)

func bigSlowOperaton() {
	defer trace("bigSlowOperaton")()
	// ...lots of work...
	time.Sleep(10 * time.Second)
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() { log.Printf("exit %s (%s),", msg, time.Since(start)) } // closure
}

func main() {
	bigSlowOperaton()
}
