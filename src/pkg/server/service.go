package server

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Service 提供服务的接口定义
type Service interface {

	// ServiceDesc 对该服务进行一些必要描述
	ServiceDesc(desc *ServiceDesc)

	// Initialization 初始化该服务所需的内容
	Initialization(ctx context.Context, state *State) error
}

// HandleFunc 服务处理函数定义
type HandleFunc func(ctx Context) (response Response)

// TwistFunc 对 gin.Context 进行扭曲变种的处理函数
type TwistFunc func(handleFunc HandleFunc) gin.HandlerFunc

// Exec 开始执行扭曲（等同直接运行该函数）
func (slf TwistFunc) Exec(handleFunc HandleFunc) gin.HandlerFunc {
	return slf(handleFunc)
}
