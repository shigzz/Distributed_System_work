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
	wservers   map[string]int
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
	wdic := make(map[string]int, len(serverList))
	for _, l := range serverList {
		wdic[l], err = GetValue(conn, l)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	fmt.Println(wdic)
	return &ClientSelector{zkaddr, serverList, wdic, selectmode}
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

var (
	wr = 0
	wi = 0
)

func (c *ClientSelector) getServerByWR() (server string, err error) {
	server = c.servers[wr]
	wi = wi + 1
	if wi >= c.wservers[server] {
		wr = (wr + 1) % len(c.servers)
		wi = 0
	}
	return
}
