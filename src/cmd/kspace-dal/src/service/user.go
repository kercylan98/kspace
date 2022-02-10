package service

import (
	"context"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/pkg/ent/ent"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/pkg/ent/entcli"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/pkg/rpc/dal"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/pkg/rpc/dal/mixin"
	"github.com/kercylan98/kspace/src/pkg/server"
	"github.com/kercylan98/kspace/src/pkg/utils/rpctime"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type User struct {
	dal.UnimplementedUserServiceServer
	MySQL entcli.MySQL
}

func (slf User) ServiceDesc(desc *server.ServiceDesc) {
	desc.RPC(&dal.UserService_ServiceDesc)
}

func (slf User) Initialization(ctx context.Context, state *server.State) error {
	slf.MySQL.Init("root:root@tcp(127.0.0.1:3306)/kspace?charset=utf8mb4&parseTime=True&loc=Local").Check()
	return slf.MySQL.Schema.Create(ctx)
}

func (slf User) Create(ctx context.Context, request *dal.CreateUserRequest) (*dal.CreateUserReply, error) {
	var userCreates = make([]*ent.UserCreate, len(request.Users))
	for i := 0; i < len(request.Users); i++ {
		userCreates = append(userCreates, slf.MySQL.User.Create().
			SetAccount(request.Users[i].Account).
			SetPassword(request.Users[i].Password))
	}
	users, err := slf.MySQL.User.CreateBulk().Save(ctx)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	var reply = new(dal.CreateUserReply)
	for i := 0; i < len(users); i++ {
		u := users[i]
		reply.Users = append(reply.Users, &dal.User{
			Id: uint32(u.ID),
			Time: &mixin.Time{
				CreatedAt: rpctime.ToTimestamp(u.CreatedAt),
				UpdatedAt: rpctime.ToTimestamp(u.UpdatedAt),
				DeletedAt: rpctime.ToTimestamp(u.DeletedAt),
			},
			Account:  u.Account,
			Password: u.Password,
		})
	}
	return reply, nil
}

func (slf User) Get(ctx context.Context, request *dal.GetUserRequest) (*dal.User, error) {
	return nil, nil
}

func (slf User) Update(ctx context.Context, request *dal.UpdateUserRequest) (*dal.User, error) {
	//TODO implement me
	panic("implement me")
}

func (slf User) Delete(ctx context.Context, request *dal.DeleteUserRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
