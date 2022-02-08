package orm

import (
	"github.com/go-redis/redis/v8"
)

// Redis 进行封装后的 Redis 客户端，使用前需要先执行"Redis.InitUse"
type Redis struct {
	*redis.Client
}

// InitUse 根据相应信息初始化默认配置的 Redis 客户端
func (slf *Redis) InitUse(addr string, password string, db int) Redis {
	slf.Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return *slf
}
