package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type client chan<- string // outgoing message channel，只负责写入到 cli
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				// 模拟耗时操作
				time.Sleep(1000 * time.Millisecond)
				cli <- msg
			}
		// 注册连入的客户端
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

// conn are shared
func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()

	// 显示在连进来的各自的客户端，提示自己的身份信息
	ch <- "You are " + who

	// 广播到所有连接的客户端用
	// 放在 entering <- ch 前，客户端第一次连接过来时，自己的终端不会显示 "... has arrived"
	messages <- who + " has arrived"

	// 注册连进来的客户端，是广播的对象，广播时消息会用到
	// 发送 ch（是 channel）
	entering <- ch

	input := bufio.NewScanner(conn)
	// Note: ignoring potential errors from input.Err()
	for input.Scan() {
		// 只会显示在自己以外的其它客户端，自己的显示内容是手动输入的
		messages <- who + ": " + input.Text()
	}

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

// 写到连进来的客户端
func clientWriter(conn net.Conn, ch <-chan string) {
	// broadcaster 里 close(cli) 后 for-range 结束，协程结束
	for msg := range ch {
		fmt.Fprintln(conn, msg) // Note: ignoring network errors
	}
}

// go build exercises-the_go_programming_language/ch8/netcat3
// go run main.go
// ./netcat3
// ./netcat3
