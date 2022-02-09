package distributed

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	basicPath = "/api"
)

// Server 基于 Zookeeper 实现的分布式服务器
type Server struct {
	Zookeeper Zookeeper
}

// NodeService 获取节点服务
func (slf Server) NodeService() NodeService {
	return NodeService{zookeeper: slf.Zookeeper}
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
