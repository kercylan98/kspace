package oauth2

import "github.com/kercylan98/kspace/src/pkg/orm"

// TokenStorage 数据存储类型约束
type TokenStorage interface {
	Redis | MySQL
}

type MySQL struct {
	orm.MySQL
	tableName  string
	gcInterval int
}

type Redis struct {
	orm.Redis
}
