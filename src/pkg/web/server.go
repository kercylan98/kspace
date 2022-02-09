package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kercylan98/kspace/src/pkg/distributed"
	"github.com/kercylan98/kspace/src/pkg/utils/netutils"
	"github.com/kercylan98/kspace/src/pkg/web/internal/utils"
	"log"
)

// New 构建并返回一个 Web 服务器
func New() Server {
	var server = Server{
		engine:      gin.Default(),
		inUseRouter: make(map[string]gin.IRouter),
	}

	server.globalErrorHandleFunc = func(ctx Context, hasErrResponse Response) {
		// By default, there is no additional processing
	}

	server.twistHandle = func(handleFunc HandleFunc) gin.HandlerFunc {
		return func(context *gin.Context) {
			ctx := Context{
				Context: context,
			}
			response := handleFunc(ctx)
			status := ctx.Writer.Status()
			if status >= 300 && status <= 399 {
				ctx.Redirect(status, ctx.Writer.Header().Get("Location"))
			} else {
				if ctx.Writer.Size() == -1 && response.noWriter == false {
					if response.Error != nil {
						response.Error.Code = status
						response.Error.Route = ctx.Request.RequestURI
						ctx.JSON(200, response)
						server.globalErrorHandleFunc(ctx, response)
					} else {
						ctx.JSON(200, response)
					}
				}
			}
		}
	}

	return server
}

// Server WEB 服务器
type Server struct {
	engine                *gin.Engine
	twistHandle           TwistFunc
	inUseRouter           map[string]gin.IRouter
	globalErrorHandleFunc func(ctx Context, hasErrResponse Response)
	node                  distributed.Node
	distributedServer     distributed.Server
	services              []Service
	routeBindFunc         []func() // 路由绑定函数
}

// DistributedRun 以分布式的方式运行服务器
func (slf *Server) DistributedRun(node distributed.Node, zookeeperHost ...string) error {
	slf.node = node
	if slf.node.IsAutoGetAddress {
		ip, err := netutils.GetOutBoundIP()
		if err != nil {
			return err
		}
		(&slf.node).Address = ip
	}
	if slf.node.IsRandomUsePort {
		port, err := netutils.GetAvailablePort()
		if err == nil {
			(&slf.node).Port = port
		}
	}

	slf.distributedServer.Zookeeper.InitUse(zookeeperHost...)
	if slf.distributedServer.Zookeeper.InitError != nil {
		return slf.distributedServer.Zookeeper.InitError
	}

	if err := slf.distributedServer.Release(slf.node); err != nil {
		return err
	}
	defer func() {
		slf.distributedServer.Close()
		log.Println("cancel release:")
	}()
	return slf.Run(fmt.Sprintf(":%d", slf.node.Port))
}

// RegisterGlobalErrorHandler 注册全局错误处理器，调用到此处之前已经对客户端进行了响应
func (slf *Server) RegisterGlobalErrorHandler(handleFunc func(ctx Context, hasErrResponse Response)) Server {
	slf.globalErrorHandleFunc = handleFunc
	return *slf
}

// RegisterMiddleware 注册中间件（middleware）到服务器中
func (slf *Server) RegisterMiddleware(middleware ...HandleFunc) Server {
	for _, m := range middleware {
		slf.engine.Use(func(context *gin.Context) {
			ctx := Context{
				Context: context,
			}
			m(ctx)
		})
	}
	return *slf
}

// RegisterVersionService 将需要指定特定版本(version)的服务(services)注册到特定名称(name)下的服务器路由器中
//
// 例如"server.RegisterVersionService("api", "v1", UserService)"，将注册到路由"/api/v1/..."下
func (slf *Server) RegisterVersionService(name string, version string, services ...Service) Server {
	name = utils.FormatUrlPathCharacter(name)
	version = utils.FormatUrlPathCharacter(version)
	routerName := fmt.Sprintf("%s:#", name)

	var router gin.IRouter
	var routerIsExist bool
	if router, routerIsExist = slf.inUseRouter[routerName]; !routerIsExist {
		router = slf.engine.Group(name).Group(version)
		slf.inUseRouter[routerName] = router
	} else {
		// name 路由器存在的情况下查找 version 路由器是否存在
		routerName = fmt.Sprintf("%s:%s", name, version)
		if router, routerIsExist = slf.inUseRouter[routerName]; !routerIsExist {
			router = slf.engine.Group(name).Group(version)
			slf.inUseRouter[routerName] = router
		}
	}

	for i := 0; i < len(services); i++ {
		service := services[i]
		if err := service.Initialization(); err != nil {
			panic(err)
		}
		slf.routeBindFunc = append(slf.routeBindFunc, func() {
			service.BindRoute(router, slf.twistHandle)
		})
		slf.services = append(slf.services, service)
	}
	fmt.Println(services)
	return *slf
}

// RegisterService 将服务(services)注册到特定名称(name)下的服务器路由器中
//
// 例如"server.RegisterService("auth", AuthService)"，将注册到路由"/auth/..."下
func (slf *Server) RegisterService(name string, services ...Service) Server {
	name = utils.FormatUrlPathCharacter(name)
	routerName := fmt.Sprintf("%s:#", name)

	var router gin.IRouter
	var routerIsExist bool
	if router, routerIsExist = slf.inUseRouter[routerName]; !routerIsExist {
		router = slf.engine.Group(name)
		slf.inUseRouter[routerName] = router
	}
	for _, service := range services {
		service.BindRoute(router, slf.twistHandle)
	}
	return *slf
}

// RegisterRootService 将服务(services)注册到根("/")路由器中
func (slf *Server) RegisterRootService(services ...Service) Server {
	for _, service := range services {
		service.BindRoute(slf.engine, slf.twistHandle)
	}
	return *slf
}

// Run 运行服务器
func (slf Server) Run(addr ...string) error {
	for i := 0; i < len(slf.services); i++ {
		slf.services[i].Runtime(Runtime{
			NodeService: slf.distributedServer.NodeService(),
		})
	}
	for _, f := range slf.routeBindFunc {
		f()
	}
	return slf.engine.Run(addr...)
}
