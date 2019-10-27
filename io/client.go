package main

import (
	"fmt"
	"net"
	"strconv"
)

func request(str string) {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("error in dialing, ", err.Error())
	}
	defer conn.Close()
	_, err = conn.Write([]byte(str))
	if err != nil {
		fmt.Println("error in writing, ", err.Error())
		return
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("err reading, ", err.Error())
	}
	fmt.Println(string(buf[:n]))
}

func main() {

	for i := 0; i < 1000; i++ {
		request("hello" + strconv.Itoa(i))
		//<-time.After(time.Second * 1)
	}
}
