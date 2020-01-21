package main

import "fmt"

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p))
	return len(p), nil
}

func main() {
	var c ByteCounter

	c.Write([]byte("hello"))
	fmt.Println(c) // 5, len("hello")

	c = 0
	var name = "Dolly"
	// Fprintf 将输入格式化后写入第一个参数（`io.Writer` interface type）
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c) // 12, len("hello, Dolly")
}
