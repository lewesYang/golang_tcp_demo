package main

import (
	"container/list"
	"fmt"
	"net"
	"time"
)

const (
	wc         = "welcome to commd chat room"
	notice     = "如需私聊请@用户名+发送内容:"
	ip         = "127.0.0.1"
	onlineUser = "在线用户:"
	port       = 5418
	buf_size   = 128
)

type client struct {
	name string
	conn net.Conn
}

var clients *list.List

func main() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 1)
	fmt.Println("servier listen on :", ip, ":", port)
	clients = list.New()
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println(err.Error())
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	fmt.Println(conn.RemoteAddr().String(), "connected!")
	onlineNotice(conn)
	c := &client{}
	c.name = conn.RemoteAddr().String()
	c.conn = conn
	clients.PushBack(c)
	for {
		data := make([]byte, 0)
		buf := make([]byte, buf_size)
		for {
			size, err := conn.Read(buf)
			if err != nil {
				continue
			}
			data = append(data, buf[:size]...)
			if size <= buf_size {
				break
			}
		}
		fmt.Println(string(data))
		go broadcast(data, conn)
	}
}

func onlineNotice(conn net.Conn) {
	sendTitle(conn)
	for e := clients.Front(); e != nil; e = e.Next() {
		conn.Write([]byte((e.Value.(*client).name)))
		conn.Write([]byte(("\r")))
	}

}

func sendTitle(conn net.Conn) {
	conn.Write([]byte(wc))
	conn.Write([]byte((notice)))
	conn.Write([]byte((onlineUser)))
}

func broadcast(data []byte, conn net.Conn) {
	for e := clients.Front(); e != nil; e = e.Next() {
		name := e.Value.(*client).name
		if name == conn.RemoteAddr().String() {
			continue
		}
		byte_name := []byte("[" + name + "]  ")
		data = append(byte_name, data...)
		e.Value.(*client).conn.Write(data)
	}
}
