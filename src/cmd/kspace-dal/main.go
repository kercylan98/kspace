package main

import (
	server "github.com/kercylan98/kspace/src/cmd/kspace-dal/src/server"
	"github.com/kercylan98/kspace/src/pkg/distributed"
	"github.com/kercylan98/kspace/src/pkg/krpc"
)

func main() {
	var err = krpc.RunDistributed(distributed.Node{
		Name:             "KSpace-DAL",
		IsAutoGetAddress: true,
		IsRandomUsePort:  true,
	}, []string{"127.0.0.1:2181"},
		server.UserServer(),
		server.OAuth2Server())
	if err != nil {
		panic(err)
	}
}
