package server

import (
	"context"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/pkg/models"
	rpc2 "github.com/kercylan98/kspace/src/cmd/kspace-dal/src/rpc"
	"github.com/kercylan98/kspace/src/pkg/cryptography"
	"github.com/kercylan98/kspace/src/pkg/krpc"
	"github.com/kercylan98/kspace/src/pkg/orm"
	"google.golang.org/grpc"
)

func OAuth2Server() krpc.ServerDemandHandlerFunc {
	return func(rpcServer grpc.ServiceRegistrar, mysql orm.MySQL, redis orm.Redis, zookeeper orm.Zookeeper, rsa *cryptography.RSA) {
		rpc2.RegisterDalOAuth2Server(rpcServer, OAuth2{
			MySQL: mysql.InitDefault("root:root@tcp(127.0.0.1:3306)/kspace?charset=utf8mb4&parseTime=True&loc=Local"),
			Redis: redis.InitUse("127.0.0.1", "root", 15),
		})
	}
}

type OAuth2 struct {
	rpc2.UnimplementedDalOAuth2Server
	MySQL orm.MySQL
	Redis orm.Redis
}

func (slf OAuth2) GetClientWithClientID(ctx context.Context, client *rpc2.AuthClient) (*rpc2.AuthClient, error) {
	oa2c := models.OAuth2Client{ClientID: client.ClientID}
	if result := slf.MySQL.Where(&oa2c).First(&oa2c); result.Error != nil {
		return nil, result.Error
	}
	return &rpc2.AuthClient{
		Id:           uint32(oa2c.ID),
		UserID:       uint32(oa2c.UserId),
		ClientID:     oa2c.ClientID,
		ClientSecret: oa2c.ClientSecret,
		Domain:       oa2c.Domain,
	}, nil
}

func (slf OAuth2) CreateClient(ctx context.Context, client *rpc2.AuthClient) (*rpc2.AuthClient, error) {
	var c = models.OAuth2Client{
		UserId:       uint(client.UserID),
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
		Domain:       client.Domain,
	}
	if result := slf.MySQL.Create(&c); result.Error != nil {
		return nil, result.Error
	}

	client.Id = uint32(c.ID)
	return client, nil
}
