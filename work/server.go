package work

import (
	"fmt"
	"io"
	"net"
	"reflect"
)

//Server 定义服务结构体
type Server struct {
	addr     string
	zkServer []string
	conn     net.Conn
	funcs    map[string]reflect.Value
}

//NewServer 创建新的Server
func NewServer(addr string, zkserver []string) *Server {
	return &Server{addr, zkserver, nil, make(map[string]reflect.Value)}
}

//Register 在服务器中注册一个rpc服务
func (s *Server) Register(name string, f interface{}) {
	if _, ok := s.funcs[name]; ok {
		//函数已经存在
		return
	}
	s.funcs[name] = reflect.ValueOf(f)
}

//Run 在zk中注册服务，并启动服务器
func (s *Server) Run() {
	var err error
	//设置监听器
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Printf("监听%s过程中发生错误: %v\n", s.addr, err)
		return
	}
	defer listener.Close()
	conn, err := GetConnect(s.zkServer)
	if err != nil {
		fmt.Printf(" connect zk error: %s ", err)
		panic(err)
	}
	err = RegistServer(conn, s.addr)
	if err != nil {
		fmt.Printf(" regist node error: %s ", err)
	}
	for {
		//接收传入的消息
		s.conn, err = listener.Accept()
		if err != nil {
			fmt.Printf("接收过程发生错误: %v\n", err)
			continue
		}
		go s.serve()
	}
}

//为每一个客户端请求分配一个服务协程
func (s *Server) serve() {
	trans := NewTransport(s.conn)

	for {
		//接收来自客户端的请求
		req, err := trans.Receive()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("读取缓存错误: %v\n", err)
			}
			return
		}
		//判断请求的函数是否存在
		f, ok := s.funcs[req.Name]
		if !ok {
			//请求的函数不存在
			e := fmt.Sprintf("请求的方法:%s不存在\n", req.Name)
			fmt.Println(e)
			continue
		}
		fmt.Printf("已请求%s函数,来自%s\n", req.Name, s.conn.RemoteAddr().String())
		//设置参数
		inArgs := make([]reflect.Value, len(req.Args))
		for i, a := range req.Args {
			inArgs[i] = reflect.ValueOf(a)
		}
		//调用真实的函数
		out := f.Call(inArgs)
		//out是来自函数的调用，最后一个参数是error，需要忽略
		outArgs := make([]interface{}, len(out)-1)
		for i := 0; i < len(out)-1; i++ {
			outArgs[i] = out[i].Interface()
		}
		var e string
		if _, ok := out[len(out)-1].Interface().(error); !ok {
			e = ""
		} else {
			e = out[len(out)-1].Interface().(error).Error()
		}
		//给客户端传送数据
		err = trans.Send(Data{Name: req.Name, Args: outArgs, Err: e})
		if err != nil {
			fmt.Printf("error in transporting data: %v\n", err)
		}
	}
}

func (s *Server) stopServer() {

}
