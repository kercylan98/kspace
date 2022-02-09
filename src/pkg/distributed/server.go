package distributed

import (
	"encoding/json"
	"fmt"
	"github.com/kercylan98/kspace/src/pkg/orm"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	basicPath = "/api"
)

// Server 基于 Zookeeper 实现的分布式服务器
type Server struct {
	Zookeeper orm.Zookeeper
}

// FindNodes 根据节点名称（name）查找已发布的节点
func (slf Server) FindNodes(name string) (nodes []Node, err error) {
	children, _, err := slf.Zookeeper.Children(fmt.Sprintf("%s/%s", basicPath, name))
	if err != nil {
		return nodes, err
	}

	for i := 0; i < i; i++ {
		child := children[i]
		data, _, err := slf.Zookeeper.Get(fmt.Sprintf("%s/%s/%s", basicPath, name, child))
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

// Release 发布项目（node）使得其可被发现
func (slf Server) Release(node Node) error {
	var (
		err        error
		serverPath = basicPath
	)
	// 检查 Zookeeper 客户端完整性
	if slf.Zookeeper.InitError != nil {
		return err
	}

	// 建立根节点
	if _, err = slf.Zookeeper.Create(serverPath, []byte(""), 0, zk.WorldACL(zk.PermAll)); err != nil && err != zk.ErrNodeExists {
		return err
	}

	// 建立服务名称节点
	serverPath = fmt.Sprintf("%s/%s", serverPath, node.Name)
	if _, err = slf.Zookeeper.Create(serverPath, []byte(""), 0, zk.WorldACL(zk.PermAll)); err != nil && err != zk.ErrNodeExists {
		return err
	}

	// 建立发布节点
	var nodeData []byte
	if nodeData, err = json.Marshal(&node); err != nil {
		return nil
	}
	serverPath = fmt.Sprintf("%s/nodes", serverPath)
	if _, err = slf.Zookeeper.CreateProtectedEphemeralSequential(serverPath, nodeData, zk.WorldACL(zk.PermAll)); err != nil && err != zk.ErrNodeExists {
		return err
	}

	return nil
}

// Close 关闭服务器
//
// 注意：需要在进程退出前调用 Server.Close 方法，否则 Zookeeper 的会话不会立即关闭，服务器创
// 建的临时节点也就不会立即消失，而是要等到timeout之后服务器才会清理。
func (slf Server) Close() {
	slf.Zookeeper.Close()
}
