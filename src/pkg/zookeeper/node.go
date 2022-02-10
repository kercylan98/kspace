package zookeeper

// Node 分布式节点信息
type Node struct {
	Name string `json:"name"` // 节点名称
	Host string `json:"host"` // 节点主机
	Port int    `json:"port"` // 节点端口
}
