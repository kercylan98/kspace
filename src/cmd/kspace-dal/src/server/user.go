package server

import (
	"context"
	"errors"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/pkg/models"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/rpc"
	"github.com/kercylan98/kspace/src/pkg/cryptography"
	"github.com/kercylan98/kspace/src/pkg/krpc"
	"github.com/kercylan98/kspace/src/pkg/model"
	"github.com/kercylan98/kspace/src/pkg/orm"
	"google.golang.org/grpc"
)

func UserServer() krpc.ServerDemandHandlerFunc {
	return func(rpcServer grpc.ServiceRegistrar, mysql orm.MySQL, redis orm.Redis, rsa *cryptography.RSA) {
		rpc.RegisterDalUserServer(rpcServer, User{
			MySQL: mysql.InitDefault("root:root@tcp(127.0.0.1:3306)/kspace?charset=utf8mb4&parseTime=True&loc=Local",
				&models.User{}, &models.OAuth2Client{},
			),
			Redis: redis.InitUse("127.0.0.1:6379", "root", 0),
			RSA:   rsa,
		})
	}
}

type User struct {
	rpc.UnimplementedDalUserServer
	MySQL orm.MySQL
	Redis orm.Redis
	RSA   *cryptography.RSA
}

func (slf User) VerifyPassword(ctx context.Context, user *rpc.User) (*rpc.User, error) {
	var replyUser, err = slf.Get(ctx, &rpc.User{
		Id:      user.Id,
		Account: user.Account,
	})
	if err != nil {
		return nil, err
	}
	if pwdData, err := slf.RSA.PrivateDecrypt(replyUser.Password); err != nil {
		return nil, err
	} else {
		if string(pwdData) == user.Password {
			return replyUser, nil
		}
	}
	return nil, errors.New("username or password incorrect")
}

func (slf User) Get(ctx context.Context, user *rpc.User) (*rpc.User, error) {
	var u = models.User{
		Core: model.Core{
			ID: uint(user.Id),
		},
		Account:       user.Account,
		OAuth2Clients: nil,
		Token:         "",
	}
	if result := slf.MySQL.Where(&u).First(&u); result.Error != nil {
		return nil, result.Error
	}

	user.Id = uint32(u.ID)
	user.Account = u.Account
	user.Password = u.Password
	return user, nil
}

func (slf User) Create(ctx context.Context, user *rpc.User) (reply *rpc.CreateReply, err error) {
	if user.Password, err = slf.RSA.PublicEncrypt([]byte(user.Password)); err != nil {
		return nil, err
	}
	u := models.User{
		Account:  user.Account,
		Password: user.Password,
	}
	if result := slf.MySQL.Create(&u); result.Error != nil {
		return nil, result.Error
	}

	user.Id = uint32(u.ID)
	return &rpc.CreateReply{
		User: user,
	}, nil
}

func (slf User) CreateMultiple(ctx context.Context, user *rpc.MultipleUser) (*rpc.CreateMultipleReply, error) {
	var err error
	var users = make([]models.User, len(user.Users))
	for i := 0; i < len(user.Users); i++ {
		var u = user.Users[i]
		u.Password, err = slf.RSA.PublicEncrypt([]byte(u.Password))
		if err != nil {
			return nil, err
		}

		users = append(users, models.User{
			Account:  u.Account,
			Password: u.Password,
		})
	}
	if result := slf.MySQL.Create(&users); result.Error != nil {
		return nil, result.Error
	}
	return &rpc.CreateMultipleReply{Size: int32(len(users))}, nil
}
