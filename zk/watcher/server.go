package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

//Server 定义服务器地址
type Server struct {
	host string
	port int
}

func main() {
	s1 := &Server{"localhost", 9897}
	s2 := &Server{"localhost", 9898}
	s3 := &Server{"localhost", 9899}
	//s4 := &Server{"192.168.3.72", 9999}
	go s1.starServer()
	go s2.starServer()
	go s3.starServer()
	//s4.addServer()
	//a := make(chan bool, 1)
	//<-a
	<-time.After(time.Second * 20)
	fmt.Println("stop server")
	//s1.stopServer()
	//s2.stopServer()
	//s3.stopServer()
	//s4.stopServer()
}

func (s *Server) starServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.host+":"+strconv.Itoa(s.port))
	fmt.Println(tcpAddr)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	defer listener.Close()
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
		_, err = conn.Write([]byte("hello world from localhost from" + conn.LocalAddr().String() + " ,number, " + strconv.Itoa(i)))
		i++
		if err != nil {
			fmt.Println("Write data error: ", err.Error())
		}
	}
}

func (s *Server) stopServer() {
	conn, err := GetConnect()
	if err != nil {
		fmt.Printf(" connect zk error: %s ", err)
	}
	defer conn.Close()
	server := s.host + ":" + strconv.Itoa(s.port)
	err = UnRegisterServer(conn, server)
	if err != nil {
		fmt.Println("err unregister server: ", err)
	}
}

func (s *Server) addServer() {
	conn, err := GetConnect()
	if err != nil {
		fmt.Printf(" connect zk error: %s ", err)
	}
	defer conn.Close()
	err = RegistServer(conn, s.host+":"+strconv.Itoa(s.port))
	if err != nil {
		fmt.Printf(" regist node error: %s ", err)
	}
}
