package work

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"time"
)

// Client struct
type Client struct {
	ClientSelector
	server string
	conn   net.Conn
}

// NewClient creates a new client
func NewClient(c ClientSelector) *Client {
	for len(c.servers) == 0 {
		fmt.Println("no servers,waiting")
		time.Sleep(1 * time.Second)
	}
	var thisServer string
	var err error
	switch c.selectmode {
	case 1:
		thisServer, err = c.getServerByRandom()
	case 2:
		thisServer, err = c.getServerByRR()
	case 3:
		thisServer, err = c.getServerByWR()
	default:
		err = errors.New("wrong select mode")
	}
	if err != nil {
		fmt.Println("get server err: ", err)
	}
	conn, err := net.Dial("tcp", thisServer)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		conn = nil
	}
	return &Client{c, thisServer, conn}
}

// Call transforms a function prototype into a function
func (c *Client) Call(name string, fptr interface{}) {
	container := reflect.ValueOf(fptr).Elem()

	f := func(req []reflect.Value) []reflect.Value {
		cliTransport := NewTransport(c.conn)

		errorHandler := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}
		// package request arguments
		inArgs := make([]interface{}, 0, len(req))
		for i := range req {
			inArgs = append(inArgs, req[i].Interface())
		}
		// send request to server
		err := cliTransport.Send(Data{Name: name, Args: inArgs})
		if err != nil { // local network error or encode error
			return errorHandler(err)
		}
		// receive response from server
		rsp, err := cliTransport.Receive()
		if err != nil { // local network error or decode error
			return errorHandler(err)
		}
		if rsp.Err != "" { // remote server error
			return errorHandler(errors.New(rsp.Err))
		}

		if len(rsp.Args) == 0 {
			rsp.Args = make([]interface{}, container.Type().NumOut())
		}
		// unpackage response arguments
		numOut := container.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 { // unpackage arguments (except error)
				if rsp.Args[i] == nil { // if argument is nil (gob will ignore "Zero" in transmission), set "Zero" value
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(rsp.Args[i])
				}
			} else { // unpackage error argument
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}

		return outArgs
	}

	container.Set(reflect.MakeFunc(container.Type(), f))
}

//GetConn 返回连接
func (c *Client) GetConn() net.Conn {
	return c.conn
}

//SyncCall 异步通信
func (c *Client) SyncCall(name string, fptr interface{}) {

}
