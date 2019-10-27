package main

import (
	"fmt"

	"time"

	"github.com/samuel/go-zookeeper/zk"
)

/*func ZkStateString(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d, Mzxid: %d, Ctime: %d, Mtime: %d, Version: %d, Cversion: %d, Aversion: %d, EphemeralOwner: %d, DataLength: %d, NumChildren: %d, Pzxid: %d",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}*/

//ZkStateStringFormat format string
func ZkStateStringFormat(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d\nMzxid: %d\nCtime: %d\nMtime: %d\nVersion: %d\nCversion: %d\nAversion: %d\nEphemeralOwner: %d\nDataLength: %d\nNumChildren: %d\nPzxid: %d\n",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}

func main() {
	fmt.Printf("ZKOperateTest\n")

	var hosts = []string{"192.168.3.71:2181"}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	var path = "/zk_test_go"
	var data = []byte("hello")
	var flags int32 = 0
	// permission
	var acls = zk.WorldACL(zk.PermAll)

	// create
	p, errCreate := conn.Create(path, data, flags, acls)
	if errCreate != nil {
		fmt.Println(errCreate)
		return
	}
	fmt.Println("created:", p)

	// get
	v, s, err := conn.Get(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("value of path[%s]=[%s].\n", path, v)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", ZkStateStringFormat(s))

	// exist
	exist, s, err := conn.Exists(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("path[%s] exist[%t]\n", path, exist)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", ZkStateStringFormat(s))

	// update
	var newData = []byte("zk_test_new_value")
	s, err = conn.Set(path, newData, s.Version)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("update state:\n")
	fmt.Printf("%s\n", ZkStateStringFormat(s))

	// get
	v, s, err = conn.Get(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("new value of path[%s]=[%s].\n", path, v)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", ZkStateStringFormat(s))

	// delete
	err = conn.Delete(path, s.Version)
	if err != nil {
		fmt.Println(err)
		return
	}

	// check exist
	exist, s, err = conn.Exists(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("after delete, path[%s] exist[%t]\n", path, exist)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", ZkStateStringFormat(s))
}
