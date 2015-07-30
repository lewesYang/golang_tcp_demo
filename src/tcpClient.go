package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	address  = "127.0.0.1:5418"
	buf_size = 128
)

func main() {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
	}
	go handle(conn)
	for {
		rd := bufio.NewReader(os.Stdin)
		line, _, _ := rd.ReadLine()
		conn.Write(line)
	}
}

func handle(conn net.Conn) {
	for {
		data := make([]byte, 0)
		buff := make([]byte, buf_size)
		for {
			size, err := conn.Read(buff)
			if err != nil {
				continue
			}
			data = append(data, buff[:size]...)
			if size <= buf_size {
				break
			}
		}
		fmt.Println(string(data))
	}
}
