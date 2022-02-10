package distributed

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
)

// NodeService 节点服务提供了一些节点操作功能
type NodeService struct {
	zookeeper Zookeeper
}

// Conn 获取特定名称节点的连接
func (slf NodeService) Conn(name string) *Conn {
	return &Conn{
		grpcConnect: map[string]*grpc.ClientConn{},
		name:        name,
		nodeService: slf,
	}
}

// Zookeeper 获取 Zookeeper
func (slf NodeService) Zookeeper() Zookeeper {
	return slf.zookeeper
}

// FindOne 根据节点名称（name）查找其中一个节点
func (slf NodeService) FindOne(name string) (Node, error) {
	if nodes, err := slf.FindNodes(name); err != nil {
		return Node{}, err
	} else {
		// TODO：简单的随机返回一个节点，应该实现负载均衡算法等。
		var tempStore = map[int]Node{}
		for i := 0; i < len(nodes); i++ {
			tempStore[i] = nodes[i]
		}
		for _, node := range tempStore {
			return node, nil
		}
	}
	return Node{}, zk.ErrNoNode
}

// FindNodes 根据节点名称（name）查找已发布的节点
func (slf NodeService) FindNodes(name string) (nodes []Node, err error) {
	children, _, err := slf.zookeeper.Children(fmt.Sprintf("%s/%s", basicPath, name))
	if err != nil {
		return nodes, err
	}

	for i := 0; i < len(children); i++ {
		child := children[i]
		data, _, err := slf.zookeeper.Get(fmt.Sprintf("%s/%s/%s", basicPath, name, child))
		if err != nil {
			if err == zk.ErrNoNode {
				continue
			}
			return nodes, err
		}
		var node Node
		if err = json.Unmarshal(data, &node); err != nil {
			return nodes, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
