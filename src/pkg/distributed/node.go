package distributed

// Node 分布式节点信息
type Node struct {
	Name    string `json:"name"`    // 节点名称
	Address string `json:"address"` // 节点地址
}
