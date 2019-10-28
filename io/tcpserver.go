package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	address := "localhost:9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	fmt.Println(tcpAddr)
	if err != nil {
		fmt.Println("err in resolve: ", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("err in listening: ", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
			continue
		}
		go handleCient(conn, address)
	}
}

var i = 1

func handleCient(conn net.Conn, port string) {
	fmt.Println("new client:", conn.RemoteAddr())
	for {
		buf := make([]byte, 1024)
		length, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading:", err.Error())
			}
			return
		}

		fmt.Println("Receive data from client:", string(buf[:length]))
		_, err = conn.Write([]byte("hello world" + strconv.Itoa(i)))
		i++
		if err != nil {
			fmt.Println("Write data error: ", err.Error())
		}
	}
}
