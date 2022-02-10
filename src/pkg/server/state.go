package server

import "sync"

// State 服务状态，可以安全的在多服务之间共享数据
type State struct {
	rw         sync.RWMutex
	anyValue   map[string]any
	serverName string // 服务名称
	serverHost string // 服务主机
	serverPort int    // 服务端口
}

func (slf *State) Set(key string, value any) {
	slf.rw.Lock()
	defer slf.rw.Unlock()
	slf.anyValue[key] = value
}

func (slf *State) Get(key string) any {
	slf.rw.RLock()
	defer slf.rw.RUnlock()
	return slf.anyValue[key]
}

func (slf *State) ServerName() string {
	return slf.serverName
}

func (slf *State) ServerHost() string {
	return slf.serverHost
}

func (slf *State) ServerPort() int {
	return slf.serverPort
}
