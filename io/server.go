package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("error in listenning, ", err.Error())
		return
	}
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("error accept, ", err.Error())
			continue
		}
		go handleConn(c)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		buf := make([]byte, 1024)
		length, err := c.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("error in reading, ", err.Error())
			}
			return
		}
		fmt.Println("Receive data from client:", string(buf[:length]))
		_, err = c.Write([]byte("hello world"))
		if err != nil {
			fmt.Println("error in writing, ", err.Error())
		}
	}
}
