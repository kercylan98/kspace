package krpc

import (
	"fmt"
	"github.com/kercylan98/kspace/src/pkg/constant"
	"github.com/kercylan98/kspace/src/pkg/cryptography"
	"github.com/kercylan98/kspace/src/pkg/distributed"
	"github.com/kercylan98/kspace/src/pkg/orm"
	"github.com/kercylan98/kspace/src/pkg/utils/netutils"
	"google.golang.org/grpc"
	"log"
	"net"
	"path/filepath"
)

// ServerDemandHandlerFunc 服务需求处理函数
type ServerDemandHandlerFunc func(rpcServer grpc.ServiceRegistrar,
	mysql orm.MySQL,
	redis orm.Redis,
	rsa *cryptography.RSA,
)

func RunDistributed(node distributed.Node, zookeeperHost []string, register ...ServerDemandHandlerFunc) error {
	if node.IsAutoGetAddress {
		ip, err := netutils.GetOutBoundIP()
		if err != nil {
			return err
		}
		(&node).Address = ip
	}
	if node.IsRandomUsePort {
		port, err := netutils.GetAvailablePort()
		if err == nil {
			(&node).Port = port
		}
	}

	var distributedServer distributed.Server
	distributedServer.Zookeeper.InitUse(zookeeperHost...)
	if distributedServer.Zookeeper.InitError != nil {
		return distributedServer.Zookeeper.InitError
	}

	if err := distributedServer.Release(node); err != nil {
		return err
	}
	defer func() {
		distributedServer.Close()
		log.Println("cancel release:", node)
	}()
	return RunServer(node.Port, register...)
}

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
		zookeeper := distributed.Zookeeper{}
		f(server, mysql, redis, rsa)
		if mysql.InitError != nil || zookeeper.InitError != nil {
			return err
		}
	}

	if lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return err
	}

	return server.Serve(lis)
}
