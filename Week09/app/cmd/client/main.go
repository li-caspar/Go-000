package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Printf("net.Dial error:%s", err)
		os.Exit(1)
	}
	defer conn.Close()
	inputReader := bufio.NewReader(os.Stdin)
	go readConn(conn)
	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("read form console faild, err:%s", err)
			break
		}
		fmt.Printf("send msg:%s", input)
		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Printf("write failed, err:%s", err)
			break
		}
	}
}

func readConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(msg)
	}
}
