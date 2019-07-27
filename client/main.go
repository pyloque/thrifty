package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pyloque/thrifty/purekv/proto/purekv"
	"strings"
)

type PureKVClient struct {
	*purekv.PureServiceClient
}

func NewPureKVClient(addr string) (*PureKVClient, error) {
	socket, _ := thrift.NewTSocket(addr)
	if err := socket.Open(); err != nil {
		panic(err)
	}
	client := purekv.NewPureServiceClientFactory(socket, thrift.NewTBinaryProtocolFactoryDefault())
	return &PureKVClient{client}, nil
}

func main() {
	addr := flag.String("addr", "localhost:8888", "remote server address")
	cmd := flag.String("cmd", "", "command with arguments")
	flag.Parse()
	client, err := NewPureKVClient(*addr)
	if err != nil {
		panic(err)
	}
	res := ""
	s := strings.SplitN(*cmd, " ", 3)
	op := strings.ToLower(s[0])
	key := s[1]
	ctx := context.Background()
	switch op {
	case "get":
		res, err = client.Get(ctx, key)
	case "set":
		value := s[2]
		res, err = client.Set(ctx, key, value)
	case "delete":
		res, err = client.Delete(ctx, key)
	default:
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("out:", res)
}

