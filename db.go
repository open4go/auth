package auth

import (
	r2redis "github.com/open4go/db/redis"
	"github.com/open4go/log"
	v9 "github.com/redis/go-redis/v9"
)

// GetRedisAuthHandler 获取数据库handler 这里定义一个方法
func GetRedisAuthHandler() *v9.Client {
	handler, err := r2redis.DBPool.GetHandler("auth")
	if err != nil {
		log.Log().Fatal(err)
	}
	return handler
}
