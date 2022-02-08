package main

import (
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/rpc"
	"github.com/kercylan98/kspace/src/cmd/kspace-uas/src/pkg/oauth2"
	"github.com/kercylan98/kspace/src/cmd/kspace-uas/src/services"
	"github.com/kercylan98/kspace/src/pkg/orm"
	"github.com/kercylan98/kspace/src/pkg/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	server := web.New()

	conn, err := grpc.Dial("127.0.0.1:9500", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	server.RegisterVersionService("api", "v1",
		createOAuth2Service(conn),
		services.Behavior{
			DalUserClient: rpc.NewDalUserClient(conn),
		},
	)

	if err := server.Run(":9501"); err != nil {
		panic(err)
	}
}

func createOAuth2Service(conn *grpc.ClientConn) services.OAuth2 {
	var (
		oauth2Service = services.OAuth2{
			DalOAuth2Client: rpc.NewDalOAuth2Client(conn),
			DalUserClient:   rpc.NewDalUserClient(conn),
		}
		oauth2Server, err = oauth2.New[oauth2.Redis](
			oauth2.Redis{Redis: (&(orm.Redis{})).InitUse("127.0.0.1:6379", "root", 15)},
			oauth2Service)
	)

	if err != nil {
		panic(err)
	}

	(&oauth2Service).Server = oauth2Server
	return oauth2Service
}
