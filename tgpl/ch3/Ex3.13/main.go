package main

import "fmt"

const (
	KB = 1 << (10 * (iota + 1))
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func main() {
	fmt.Println(KB, MB, EB)
}
