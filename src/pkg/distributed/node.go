package distributed

// Node 分布式节点信息
type Node struct {
	Name    string `json:"name"`    // 节点名称
	Address string `json:"address"` // 节点地址
	Port    int    `json:"port"`    // 节点端口

	IsAutoGetAddress bool `json:"-"` // 是否自动获取发布的IP地址
	IsRandomUsePort  bool `json:"-"` // 是否随机使用可用端口
}
