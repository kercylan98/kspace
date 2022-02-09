package web

import "github.com/kercylan98/kspace/src/pkg/distributed"

// Runtime 服务运行时状态
type Runtime struct {
	NodeService distributed.NodeService
}
