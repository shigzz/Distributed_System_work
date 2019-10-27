package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

/*func ZkStateString(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d, Mzxid: %d, Ctime: %d, Mtime: %d, Version: %d, Cversion: %d, Aversion: %d, EphemeralOwner: %d, DataLength: %d, NumChildren: %d, Pzxid: %d",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}*/

//TestChildren func
func TestChildren(t *testing.T) {
	fmt.Printf("ZkChildWatchTest")

	var hosts = []string{"192.168.0.72:2181"}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// try create root path
	var rootPath = "/test_root"

	// check root path exist
	exist, _, err := conn.Exists(rootPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !exist {
		fmt.Printf("try create root path: %s\n", rootPath)
		var acls = zk.WorldACL(zk.PermAll)
		p, err := conn.Create(rootPath, []byte("root_value"), 0, acls)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("root_path: %s create\n", p)
	}

	// try create child node
	curTime := time.Now().Unix()
	chPath := fmt.Sprintf("%s/ch_%d", rootPath, curTime)
	var acls = zk.WorldACL(zk.PermAll)
	p, err := conn.Create(chPath, []byte("ch_value"), zk.FlagEphemeral, acls)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("ch_path: %s create\n", p)

	// watch the child events
	children, s, childCh, err := conn.ChildrenW(rootPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("root_path[%s] child_count[%d]\n", rootPath, len(children))
	for idx, ch := range children {
		fmt.Printf("%d, %s\n", idx, ch)
	}

	fmt.Printf("watch children result state[%s]\n", ZkStateString(s))

	for {
		select {
		case chEvent := <-childCh:
			{
				fmt.Println("path:", chEvent.Path)
				fmt.Println("type:", chEvent.Type.String())
				fmt.Println("state:", chEvent.State.String())

				if chEvent.Type == zk.EventNodeCreated {
					fmt.Printf("has node[%s] detete\n", chEvent.Path)
				} else if chEvent.Type == zk.EventNodeDeleted {
					fmt.Printf("has new node[%d] create\n", chEvent.Path)
				} else if chEvent.Type == zk.EventNodeDataChanged {
					fmt.Printf("has node[%d] data changed", chEvent.Path)
				}
			}
		}

		time.Sleep(time.Millisecond * 10)
	}
}
