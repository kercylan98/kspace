package main

import (
	"github.com/kercylan98/kspace/src/cmd/kspace-uas/src/services"
	"github.com/kercylan98/kspace/src/pkg/web"
)

func main() {
	server := web.New()

	server.RegisterVersionService("api", "v1",
		new(services.OAuth2),
		new(services.Behavior),
	)

	if err := server.Run(":9501"); err != nil {
		panic(err)
	}
}
