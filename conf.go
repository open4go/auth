package auth

import (
	"os"
	"time"
)

const (
	// TokenExpireTimeConfKey 默认配置
	TokenExpireTimeConfKey = "TOKEN_EXPIRE_TIME_CONF_KEY"
	defaultExpireTime      = time.Minute * 30
)

func getExpireTime() time.Duration {
	t := os.Getenv(TokenExpireTimeConfKey)
	if t != "" {
		val, err := time.ParseDuration(t)
		if err != nil {
			return defaultExpireTime
		}
		return val
	}
	return defaultExpireTime
}
