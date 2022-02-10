package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kercylan98/kspace/src/pkg/utils/netutils"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

// Server 针对 WEB、RPC等服务的统一实现
type Server struct {
	initLock                  sync.Once                                    // 鉴定 Server 是否已经初始化的锁
	rpc                       *grpc.Server                                 // 基于 grpc.Server 的 RPC 服务器
	web                       *gin.Engine                                  // 基于 gin.Engine 的 WEB 服务器
	webRouteInfo              map[Service][]ServiceDescWebRouteHandlerFunc // 记录了 WEB 服务路由信息，将在初始化后调用并清理
	services                  []Service                                    // 所有被采用的服务集合（不区分类别）
	twist                     TwistFunc                                    // 对于 gin.Context 上下文进行扭曲的函数
	state                     *State                                       // 服务状态
	discovery                 Discovery                                    // 服务发现实现
	discoveryErrorHandlerFunc []func(err error)                            // 服务发现错误处理函数
}

// init 初始化服务器
func (slf *Server) init() {
	slf.webRouteInfo = make(map[Service][]ServiceDescWebRouteHandlerFunc)
	slf.state = &State{
		anyValue: make(map[string]any),
	}
	slf.twist = func(handleFunc HandleFunc) gin.HandlerFunc {
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
					} else {
						ctx.JSON(200, response)
					}
				}
			}
		}
	}

}

// Use 定义服务器所使用的服务
func (slf *Server) Use(services ...Service) *Server {
	slf.initLock.Do(slf.init)

	for _, service := range services {
		var serviceDesc ServiceDesc
		service.ServiceDesc(&serviceDesc)

		if len(serviceDesc.grpcServiceDesc) > 0 && slf.rpc == nil {
			slf.rpc = grpc.NewServer()
		}
		for _, desc := range serviceDesc.grpcServiceDesc {
			slf.rpc.RegisterService(desc, service)
		}

		if len(serviceDesc.routeBindHandlerFuncList) > 0 {
			slf.webRouteInfo[service] = serviceDesc.routeBindHandlerFuncList
		}

		slf.services = append(slf.services, service)
	}
	return slf
}

// AddState 添加服务器状态
func (slf *Server) AddState(key string, value any) *Server {
	slf.state.Set(key, value)
	return slf
}

// Discovery 使得该服务器允许被服务治理发现
func (slf *Server) Discovery(name string, discovery Discovery, discoveryErrorHandlerFunc ...func(err error)) {
	slf.state.serverName = name
	slf.discovery = discovery
	slf.discoveryErrorHandlerFunc = discoveryErrorHandlerFunc
}

// Run 以指定端口（port）运行服务器，如果没有任何端口，将会随机选择一个可用端口运行，如果多个端口，仅首个生效
func (slf *Server) Run(port ...int) error {
	var ctx = context.Background()
	if len(slf.webRouteInfo) > 0 && slf.web == nil {
		slf.web = gin.Default()
	}
	for _, service := range slf.services {
		if err := service.Initialization(ctx, slf.state); err != nil {
			return err
		}
		if routeInfo, exist := slf.webRouteInfo[service]; exist {
			for _, handlerFunc := range routeInfo {
				handlerFunc(slf.web, slf.twist)
			}
		}
	}
	slf.webRouteInfo = nil

	var p int
	if len(port) > 0 {
		p = port[0]
	} else {
		ap, err := netutils.GetAvailablePort()
		if err != nil {
			return err
		}
		p = ap
	}

	ip, err := netutils.GetOutBoundIP()
	if err != nil {
		return err
	}
	slf.state.serverHost = ip
	slf.state.serverPort = p

	listen, err := net.Listen("tcp", fmt.Sprint(":", p))
	if err != nil {
		return err
	}

	if slf.discovery != nil {
		if errChan, err := slf.discovery.Release(ctx, slf.state); err != nil {
			return err
		} else {
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						if err := <-errChan; len(slf.discoveryErrorHandlerFunc) > 0 {
							for i := 0; i < len(slf.discoveryErrorHandlerFunc); i++ {
								slf.discoveryErrorHandlerFunc[i](err)
							}
						} else {
							log.Println("[WARN] use server discovery, but no use err handler. err:", err)
						}
					}
				}
			}()
		}
	}

	mux := cmux.New(listen)

	if slf.rpc != nil {
		rpc := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
		go func() {
		run:
			{
				if err := slf.rpc.Serve(rpc); err != nil {
					goto run
				}
			}
		}()
	}

	if slf.web != nil {
		web := mux.Match(cmux.Any())
		go func() {
		run:
			{
				if err := slf.web.RunListener(web); err != nil {
					goto run
				}
			}
		}()
	}

	return mux.Serve()
}
