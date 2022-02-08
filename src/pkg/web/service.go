package web

import (
	"github.com/gin-gonic/gin"
)

// Service 提供 Web 服务的接口定义
type Service interface {

	// Initialization 初始化该服务所需的内容（将在 BindRoute 之前进行）
	Initialization() error

	// BindRoute 定义该服务需要绑定的路由
	BindRoute(router gin.IRouter, twist TwistFunc)
}

// HandleFunc 服务处理函数定义
type HandleFunc func(ctx Context) (response Response)

// TwistFunc 对 gin.Context 进行扭曲变种的处理函数
type TwistFunc func(handleFunc HandleFunc) gin.HandlerFunc

// Exec 开始执行扭曲（等同直接运行该函数）
func (slf TwistFunc) Exec(handleFunc HandleFunc) gin.HandlerFunc {
	return slf(handleFunc)
}
