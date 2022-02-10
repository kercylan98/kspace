package entcli

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/pkg/ent/ent"
)

// MySQL 进行封装后的 MySQL 客户端，使用前需要先执行"MySQL.InitDefault"或"MySQL.InitUse"函数
type MySQL struct {
	*ent.Client
	InitError error
}

// Init 根据 dsn 初始化默认配置的 MySQL 客户端
//
// DSN-Template: root:root@tcp(127.0.0.1:3306)/kspace?charset=utf8mb4&parseTime=True&loc=Local
func (slf *MySQL) Init(dsn string, option ...ent.Option) MySQL {
	slf.Client, slf.InitError = ent.Open("mysql", dsn, option...)
	return *slf
}

// Check 检查并返回本实例，如果存在异常将 panic
func (slf MySQL) Check() MySQL {
	if slf.InitError != nil {
		panic(slf.InitError)
	}
	return slf
}
