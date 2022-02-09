package distributed

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Conn struct {
	grpcConnect map[string]*grpc.ClientConn
	name        string
	nodeService NodeService
}

func (slf *Conn) createGRpcConn() (*grpc.ClientConn, error) {
	node, err := slf.nodeService.FindOne(slf.name)
	if err != nil {
		return nil, err
	}
	address := fmt.Sprintf("%s:%d", node.Address, node.Port)

	if gc, exist := slf.grpcConnect[address]; exist {
		return gc, nil
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	slf.grpcConnect[address] = conn
	return conn, nil
}

func (slf *Conn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	conn, err := slf.createGRpcConn()
	if err != nil {
		return err
	}
	return conn.Invoke(ctx, method, args, reply, opts...)
}

func (slf *Conn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	conn, err := slf.createGRpcConn()
	if err != nil {
		return nil, err
	}
	return conn.NewStream(ctx, desc, method, opts...)
}
