package server

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// ServiceDescWebRouteHandlerFunc 服务 WEB 路由描述处理函数
type ServiceDescWebRouteHandlerFunc func(router gin.IRouter, twist TwistFunc)

// ServiceDesc 混合 grpc.ServiceDesc 和 gin.IRouter 的服务描述结构
type ServiceDesc struct {
	grpcServiceDesc          []*grpc.ServiceDesc              // 存储了该服务实现了哪些 RPC 服务的信息
	routeBindHandlerFuncList []ServiceDescWebRouteHandlerFunc // 存储了该服务提供的 WEB 路由信息
}

// RPC 描述 RPC 服务信息
func (slf *ServiceDesc) RPC(desc ...*grpc.ServiceDesc) {
	slf.grpcServiceDesc = desc
}

// WEB 描述 WEB 服务信息
func (slf *ServiceDesc) WEB(routeHandlerFunc ...ServiceDescWebRouteHandlerFunc) {
	slf.routeBindHandlerFuncList = routeHandlerFunc
}
