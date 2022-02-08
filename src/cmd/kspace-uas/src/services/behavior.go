package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/pkg/models"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/rpc"
	"github.com/kercylan98/kspace/src/pkg/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Behavior 用户行为服务
type Behavior struct {
	rpc.DalUserClient
}

func (slf *Behavior) Initialization() error {
	var (
		err  error
		conn *grpc.ClientConn
	)

	if conn, err = grpc.Dial("127.0.0.1:9500", grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return err
	}

	slf.DalUserClient = rpc.NewDalUserClient(conn)
	return nil
}

func (slf Behavior) BindRoute(router gin.IRouter, twist web.TwistFunc) {
	router.Group("/behavior").
		POST("/signup", twist.Exec(slf.Signup))
}

// Signup 用户注册
func (slf Behavior) Signup(ctx web.Context) (response web.Response) {
	user := models.User{}
	if err := ctx.ShouldBind(&user); err != nil {
		return response.Err(err).
			MaybeSo("check whether the incoming parameters are as expected").
			Show("用户注册失败，请稍后再试。").
			Throw()
	}

	reply, err := slf.Create(ctx.Request.Context(), &rpc.User{
		Account:  user.Account,
		Password: user.Password,
	})

	if err != nil {
		return response.Err(err).
			MaybeSo("an exception occurred during database insertion, usually because the account already exists").
			Show(fmt.Sprintf("用户注册失败，该账号（%s）已存在", user.Account)).
			Throw()
	}

	(&user).ID = uint(reply.User.Id)

	return response.Pass(user)
}
