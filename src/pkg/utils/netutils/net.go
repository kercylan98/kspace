package netutils

import (
	"fmt"
	"net"
	"strings"
)

// GetAvailablePort 获取可用端口
func GetAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	return listener.Addr().(*net.TCPAddr).Port, listener.Close()

}

// IsPortAvailable 判断端口是否可以（未被占用）
func IsPortAvailable(port int) (bool, error) {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false, err
	}

	return true, listener.Close()
}

// GetOutBoundIP 获取对外IP地址
func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "", err
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
