package main

import (
	"context"
	"flag"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pyloque/thrifty/purekv/proto/purekv"
	"log"
	"sync"
)

type PureDB struct {
	sync.Mutex
	kvs map[string]string
}

func (db *PureDB) Get(ctx context.Context, key string) (r string, err error) {
	db.Lock()
	defer db.Unlock()
	return db.kvs[key], nil
}

func (db *PureDB) Set(ctx context.Context, key string, value string) (r string, err error) {
	db.Lock()
	defer db.Unlock()
	var oldValue = db.kvs[key]
	db.kvs[key] = value
	return oldValue, nil
}

func (db *PureDB) Delete(ctx context.Context, key string) (r string, err error) {
	db.Lock()
	defer db.Unlock()
	var oldValue = db.kvs[key]
	delete(db.kvs, key)
	return oldValue, nil
}

func NewPureDB() *PureDB {
	return &PureDB {
		kvs: make(map[string]string),
	}
}

func main() {
	addr := flag.String("addr", "localhost:8888", "server listen address")
	flag.Parse()
	transport, _ := thrift.NewTServerSocket(*addr)
	processor := purekv.NewPureServiceProcessor(NewPureDB())
	s := thrift.NewTSimpleServer2(processor, transport)
	err := s.Listen()
	if err!= nil {
		log.Panic("address busy")
	}
	log.Printf("server started on %s", *addr)
	_ = s.AcceptLoop()
}
