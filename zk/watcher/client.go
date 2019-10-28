package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

var serverList []string

func main() {
	conn, err := GetConnect()
	if err != nil {
		fmt.Printf(" connect zk error: %s \n ", err)
		return
	}
	defer conn.Close()
	serverList, err = GetServerList(conn)
	if err != nil {
		fmt.Printf(" get server list error: %s \n", err)
		return
	}

	/*count := len(serverList)
	if count == 0 {
		err = errors.New("server list is empty")
		return
	}*/

	//用来实时监听服务的上线与下线功能，serverList时刻保持最新的在线服务
	snapshots, errors := watchServerList(conn, "/go_servers")
	go func() {
		/*for {
			select {
			case serverList = <-snapshots:
				fmt.Printf("1111:%+v\n", serverList)
				go start()
			case err := <-errors:
				fmt.Printf("2222:%+v\n", err)
			}
		}*/
		for list := range snapshots {
			serverList = list
			fmt.Println("11111", serverList)
			//start()
		}
		//<-time.After(time.Second * 100)
	}()

	MakeDir(conn, "/config")

	configs, errors := watchGetDat(conn, "/config")
	go func() {
		for {
			select {
			case configData := <-configs:
				fmt.Printf("333:%+v\n", string(configData))
			case err := <-errors:
				fmt.Printf("4444:%+v\n", err)
			}
		}
	}()

	/*for {
		time.Sleep(1 * time.Second)
	}*/

	for len(serverList) == 0 {
		time.Sleep(1 * time.Second)
	}

	for i := 0; i < 100; i++ {
		fmt.Println("start Client :", i)

		startClient()

		time.Sleep(1 * time.Second)
	}
}

func start() {
	for i := 0; i < 10; i++ {
		fmt.Println("start Client :", i)
		startClient()
		time.Sleep(1 * time.Second)
	}
}

func startClient() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err:", err)
		}
	}()
	// service := "127.0.0.1:8899"
	//获取地址
	serverHost, err := getServerHost()
	if err != nil {
		fmt.Printf("get server host fail: %s \n", err)
		return
	}
	//serverHost := "127.0.0.1:8899"
	fmt.Println("connect host: " + serverHost)
	//tcpAddr, err := net.ResolveTCPAddr("tcp4", serverHost)
	//checkError(err)
	conn, err := net.Dial("tcp", serverHost)
	checkError(err)
	defer conn.Close()
	fmt.Println("connect ok")
	_, err = conn.Write([]byte("timestamp"))
	checkError(err)
	fmt.Println("write ok")
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	checkError(err)
	fmt.Println("recv:", string(buf[:n]))

	return
}

func getServerHost() (host string, err error) {
	//随机选中一个返回
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	host = serverList[r.Intn(3)]
	return
}
