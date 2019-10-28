package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

//Server 定义服务器地址
type Server struct {
	host string
	port int
}

func main() {
	s1 := &Server{"localhost", 9697}
	s2 := &Server{"localhost", 9698}
	s3 := &Server{"localhost", 9699}
	go s1.starServer()
	go s2.starServer()
	go s3.starServer()

	a := make(chan bool, 1)
	<-a
}

func (s *Server) starServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.host+":"+strconv.Itoa(s.port))
	fmt.Println(tcpAddr)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	//注册zk节点q
	conn, err := GetConnect()
	if err != nil {
		fmt.Printf(" connect zk error: %s ", err)
	}
	defer conn.Close()
	err = RegistServer(conn, s.host+":"+strconv.Itoa(s.port))
	if err != nil {
		fmt.Printf(" regist node error: %s ", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
			continue
		}
		go handleCient(conn, s.host+":"+strconv.Itoa(s.port))
	}

	fmt.Println("aaaaaa")
}

func handleCient(conn net.Conn, port string) {
	fmt.Println("new client:", conn.RemoteAddr())
	i := 1
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
		if err != nil {
			fmt.Println("Write data error: ", err.Error())
		}
	}
}
