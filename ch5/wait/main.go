// attempts to contact the server of URL; tries for one minute
// using exponential back-off

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func WaitForServer(url string) error {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)

	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil
		}
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}

func main() {
	if err := WaitForServer("http://hofungkoeng.com"); err != nil {
		log.Fatalf("Site is down: %v\n", err)
	}

	/* if err := WaitForSever(url); err != nil {
	    fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
	    os.Exit(1)
	} */
}
