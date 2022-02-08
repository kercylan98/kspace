package main

import (
	server2 "github.com/kercylan98/kspace/src/cmd/kspace-dal/src/server"
	"github.com/kercylan98/kspace/src/pkg/krpc"
)

func main() {
	var err = krpc.RunServer(9500,
		server2.UserServer(),
		server2.OAuth2Server())
	if err != nil {
		panic(err)
	}
}
