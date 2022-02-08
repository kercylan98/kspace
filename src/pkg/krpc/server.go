package krpc

import (
	"fmt"
	"github.com/kercylan98/kspace/src/pkg/constant"
	"github.com/kercylan98/kspace/src/pkg/cryptography"
	"github.com/kercylan98/kspace/src/pkg/orm"
	"google.golang.org/grpc"
	"net"
	"path/filepath"
)

// ServerDemandHandlerFunc 服务需求处理函数
type ServerDemandHandlerFunc func(rpcServer grpc.ServiceRegistrar,
	mysql orm.MySQL,
	redis orm.Redis,
	rsa *cryptography.RSA,
)

// RunServer 运行RPC服务器
func RunServer(port int, register ...ServerDemandHandlerFunc) error {
	var (
		err                       error
		server                    = grpc.NewServer()
		lis                       net.Listener
		rsaSecretKeyDirectoryPath string
		rsa                       *cryptography.RSA
	)

	if rsaSecretKeyDirectoryPath, err = filepath.Abs(constant.RSASecretKeyDirectoryPath); err != nil {
		return err
	}

	if rsa, err = cryptography.NewRsaWithFile(
		filepath.Join(rsaSecretKeyDirectoryPath, "public.key"),
		filepath.Join(rsaSecretKeyDirectoryPath, "private.key")); err != nil {
		return err
	}

	for _, f := range register {
		mysql := orm.MySQL{}
		redis := orm.Redis{}
		if f(server, mysql, redis, rsa); mysql.InitError != nil {
			return err
		}
	}

	if lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return err
	}

	return server.Serve(lis)
}
