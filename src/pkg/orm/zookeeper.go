package orm

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

// Zookeeper 进行封装后的 Zookeeper 客户端，使用前需要先执行"Zookeeper.InitUse"
type Zookeeper struct {
	*zk.Conn
	Event     <-chan zk.Event
	InitError error
}

// InitUse 根据相应信息初始化默认配置的 Zookeeper 客户端（需要处理 Zookeeper.InitError）
func (slf *Zookeeper) InitUse(hosts ...string) Zookeeper {
	slf.Conn, slf.Event, slf.InitError = zk.Connect(hosts, 5*time.Second)
	if slf.InitError != nil {
		return *slf
	}
	return *slf
}
