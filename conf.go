package auth

import (
	"os"
	"strconv"
	"time"
)

const (
	// TokenExpireTimeConfKey 默认配置
	TokenExpireTimeConfKey = "TOKEN_EXPIRE_TIME_CONF_KEY"
	defaultExpireTime      = time.Minute * 30
)

func getExpireTime() int {
	t := os.Getenv(TokenExpireTimeConfKey)
	if t != "" {
		val, err := strconv.Atoi(t)
		if err != nil {
			return int(defaultExpireTime)
		}
		return val
	}
	return int(defaultExpireTime)
}
