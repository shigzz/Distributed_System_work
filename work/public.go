package work

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

//GetConnect 获取连接
func GetConnect(hosts []string) (conn *zk.Conn, err error) {
	//hosts := []string{"192.168.3.71:2181"}
	conn, _, err = zk.Connect(hosts, 5*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	return
}

//MakeDir 注册路径
func MakeDir(conn *zk.Conn, path string) (err error) {
	if path == "" {
		return errors.New("path should not been empty")
	}
	if path == "/" {
		return nil
	}
	if path[0] != '/' {
		return errors.New("path must start with /")
	}

	exist, _, err := conn.Exists(path)
	if exist {
		return nil
	}
	_, err = conn.Create(path, []byte(""), int32(0), zk.WorldACL(zk.PermAll))
	if err == nil {
		return nil
	}

	//从父节点开始创建节点
	paths := strings.Split(path[1:], "/")
	createdPath := ""
	for _, p := range paths {
		createdPath = createdPath + "/" + p
		exist, _, err = conn.Exists(createdPath)
		if !exist {
			_, err = conn.Create(createdPath, []byte(""), int32(0), zk.WorldACL(zk.PermAll))
			if err != nil {
				return
			}
		}
	}
	return nil
}

//DeletePath 删除path
func DeletePath(conn *zk.Conn, path string) (err error) {
	exist, s, err := conn.Exists(path)
	if !exist {
		return nil
	}
	err = conn.Delete(path, s.Version)
	fmt.Println(err, "bbb")
	if err != nil {
		return err
	}
	return nil
}

//RegistServer 注册server
func RegistServer(conn *zk.Conn, host string, w int) (err error) {
	path := "/" + "go_servers"
	err = MakeDir(conn, path)
	if err != nil {
		return
	}
	path = path + "/" + host
	wei := strconv.Itoa(w)
	_, err = conn.Create(path, []byte(wei), int32(1), zk.WorldACL(zk.PermAll))
	return
}

//UnRegisterServer 取消注册
func UnRegisterServer(conn *zk.Conn, host string) (err error) {
	path := "/" + "go_servers" + "/" + host
	err = DeletePath(conn, path)
	fmt.Println(err, "ccc", path)
	if err != nil {
		return
	}
	return nil
}

//GetServerList 获取服务列表
func GetServerList(conn *zk.Conn) (list []string, err error) {
	list, _, err = conn.Children("/go_servers")
	fmt.Println("list:", list)
	return
}

//GetValue 获取值
func GetValue(conn *zk.Conn, p string) (int, error) {
	path := "/" + "go_servers/" + p
	v, _, err := conn.Get(path)
	if err != nil {
		return 0, err
	}
	w, err := strconv.Atoi((string(v)))
	if err != nil {
		return 0, err
	}
	return w, nil
}

//WatchServerList watch机制，服务器有断开或者重连，收到消息
func WatchServerList(conn *zk.Conn, path string) (chan []string, chan error) {
	snapshots := make(chan []string)
	errors := make(chan error)

	go func() {
		for {
			snapshot, _, events, err := conn.ChildrenW(path)
			if err != nil {
				errors <- err
				return
			}
			snapshots <- snapshot
			evt := <-events
			if evt.Err != nil {
				errors <- evt.Err
				return
			}
		}
	}()

	return snapshots, errors
}

//watch机制，监听配置文件变化的过程
func watchGetDat(conn *zk.Conn, path string) (chan []byte, chan error) {
	snapshots := make(chan []byte)
	errors := make(chan error)

	go func() {
		for {
			dataBuf, _, events, err := conn.GetW(path)
			if err != nil {
				errors <- err
				return
			}
			snapshots <- dataBuf
			evt := <-events
			if evt.Err != nil {
				errors <- evt.Err
				return
			}
		}
	}()

	return snapshots, errors
}

func checkError(err error) {
	if err != nil {
		fmt.Println("err:", err)
		//panic(err)
	}
}
