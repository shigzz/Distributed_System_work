package work

import (
	"fmt"
	"math/rand"
	"time"
)

//ClientSelector 选择服务器
type ClientSelector struct {
	zkaddr     []string
	servers    []string
	selectmode int
}

//NewClientSelector 新建selector对象
func NewClientSelector(zkaddr []string, selectmode int) *ClientSelector {
	conn, err := GetConnect(zkaddr)
	if err != nil {
		fmt.Println("failed to connect zookeeper server: ", err)
		panic(err)
	}
	defer conn.Close()
	serverList, err := GetServerList(conn)
	if err != nil {
		fmt.Printf(" get server list error: %s \n", err)
		panic(err)
	}
	//监听服务器列表的变化，以及错误
	//snapshots, errors := watchServerList(conn, "/go_servers")
	return &ClientSelector{zkaddr, serverList, selectmode}
}

func (c *ClientSelector) getServerByRandom() (server string, err error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := len(c.servers)
	server = c.servers[r.Intn(length)]
	return
}

var rr int

func (c *ClientSelector) getServerByRR() (server string, err error) {
	rr = (rr + 1) % len(c.servers)
	server = c.servers[rr]
	fmt.Println(rr, server)
	return
}
