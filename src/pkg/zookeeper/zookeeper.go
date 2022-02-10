package zookeeper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/kercylan98/kspace/src/pkg/server"
	"log"
	"sync"
	"time"
)

// Zookeeper 进行封装后的 Zookeeper 客户端，使用前需要先执行"Zookeeper.InitUse"
type Zookeeper struct {
	*zk.Conn
	Event     <-chan zk.Event
	InitError error

	nodeTable     map[string][]Node // 节点表
	nodeTableLock sync.RWMutex      // 节点表锁
}

// InitUse 根据相应信息初始化默认配置的 Zookeeper 客户端（需要处理 Zookeeper.InitError）
func (slf *Zookeeper) InitUse(hosts ...string) *Zookeeper {
	slf.Conn, slf.Event, slf.InitError = zk.Connect(hosts, 5*time.Second)
	if slf.InitError != nil {
		return slf
	}
	return slf
}

// Check 检查并返回本实例，如果存在异常将 panic
func (slf *Zookeeper) Check() *Zookeeper {
	if slf.InitError != nil {
		panic(slf.InitError)
	}
	return slf
}

// RefreshNodeTable 刷新节点缓存表
func (slf *Zookeeper) RefreshNodeTable() error {
	nodes, _, err := slf.Children("/discovery")
	if err != nil {
		return err
	}

	slf.nodeTableLock.Lock()
	defer slf.nodeTableLock.Unlock()
	slf.nodeTable = map[string][]Node{}
	for _, nodeID := range nodes {
		data, _, err := slf.Get(fmt.Sprintf("/discovery/%s", nodeID))
		if err != nil {
			return err
		}
		var node Node
		if err = json.Unmarshal(data, &node); err != nil {
			return err
		}
		if _, exist := slf.nodeTable[node.Name]; !exist {
			slf.nodeTable[node.Name] = []Node{}
		}
		slf.nodeTable[node.Name] = append(slf.nodeTable[node.Name], node)
	}
	return nil
}

// Release 发布服务器（state）使得其可被发现
func (slf *Zookeeper) Release(ctx context.Context, state *server.State) (watchErr <-chan error, err error) {
	// 检查 Zookeeper 客户端完整性
	if slf.InitError != nil {
		return nil, err
	}

	// 建立全局根节点
	if _, err = slf.Create("/discovery", []byte(""), 0, zk.WorldACL(zk.PermAll)); err != nil && err != zk.ErrNodeExists {
		return nil, err
	} else {
		log.Println("node found server discovery root node, init it")
	}

	// 建立发布节点
	var node = Node{
		Name: state.ServerName(),
		Host: state.ServerHost(),
		Port: state.ServerPort(),
	}
	var nodeData []byte
	if nodeData, err = json.Marshal(&node); err != nil {
		return nil, err
	}
	if _, err = slf.CreateProtectedEphemeralSequential("/discovery/nodes", nodeData, zk.WorldACL(zk.PermAll)); err != nil && err != zk.ErrNodeExists {
		return nil, err
	}

	// 刷新一次节点
	if err = slf.RefreshNodeTable(); err != nil {
		return nil, err
	}

	// 监听节点更新
	we := make(chan error)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, _, event, err := slf.Conn.ChildrenW("/discovery")
				if err != nil {
					we <- err
					break
				}

				evt := <-event
				if evt.Type == zk.EventNodeChildrenChanged {
					if err = slf.RefreshNodeTable(); err != nil {
						we <- err
						break
					}
				} else {
					fmt.Println(evt)
				}

			}
		}
	}()

	return we, nil
}
