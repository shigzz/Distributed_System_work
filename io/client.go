package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
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
	reschan := getResponse(conn)
	/*go func() {
		for res := range reschan {
			fmt.Println(res)
			break
		}
	}()*/
	//fmt.Println(str)
	//time.Sleep(1 * time.Second)
	/*for res := range reschan {
		fmt.Println(res)
		break
	}*/

}

func getResponse(conn net.Conn) chan string {
	buf := make([]byte, 1024)
	res := make(chan string)
	go func() {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err reading, ", err.Error())
		}
		//fmt.Println(string(buf[:n]))
		res <- string(buf[:n])
		time.Sleep(time.Second * 1)
	}()
	return res
}

func main() {

	for i := 0; i < 50; i++ {
		request("hello" + strconv.Itoa(i))
		//time.Sleep(time.Second * 2)
	}
	<-time.After(time.Second * 20)
}
