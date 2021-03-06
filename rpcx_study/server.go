package main

import "github.com/smallnest/rpcx/server"

type Args struct {
	A int `msg:"a"`
	B int `msg:"b"`
}
type Reply struct {
	C int `msg:"c"`
}
type Arith int

func (t *Arith) Mul(args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}
func (t *Arith) Error(args *Args, reply *Reply) error {
	panic("ERROR")
}
func main() {
	server := server.NewServer()
	server.RegisterName("Arith", new(Arith), "0")
	server.Serve("tcp", "127.0.0.1:8972")
}
