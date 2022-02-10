package main

import (
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/service"
	"github.com/kercylan98/kspace/src/pkg/server"
	"github.com/kercylan98/kspace/src/pkg/zookeeper"
)

func main() {
	srv := new(server.Server)
	srv.Use(new(service.User))
	srv.Discovery("KSPACE-DAL", new(zookeeper.Zookeeper).InitUse("127.0.0.1:2181").Check())
	if err := srv.Run(); err != nil {
		panic(err)
	}
}
