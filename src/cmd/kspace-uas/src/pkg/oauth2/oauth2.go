package oauth2

import (
	"github.com/go-oauth2/mysql/v4"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	redisStore "github.com/go-oauth2/redis/v4"
)

// Server 基于 OAuth2 的鉴权服务器
type Server[T TokenStorage] struct {
	*server.Server
}

// New 通过指定 Token 存储（storage）要求的 TokenStorage 类型和客户端存储（clientStore）来构建 OAuth2 服务器并返回错误信息
func New[T TokenStorage](storage T, clientStore oauth2.ClientStore) (Server[T], error) {
	manager := manage.NewDefaultManager()

	switch v := any(storage).(type) {
	case MySQL:
		if db, err := v.DB.DB(); err != nil {
			return Server[T]{}, err
		} else {
			manager.MapTokenStorage(mysql.NewStoreWithDB(db, v.tableName, v.gcInterval))
		}
	case Redis:
		manager.MapTokenStorage(redisStore.NewRedisStoreWithCli(v.Client))
	}

	manager.MapClientStorage(clientStore)
	srv := server.NewDefaultServer(manager)

	return Server[T]{Server: srv}, nil
}
