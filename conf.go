package auth

import (
	"github.com/spf13/viper"
	"time"
)

const (
	defaultExpireTime = time.Minute * 30
)

func getExpireTime() time.Duration {
	t := viper.GetString("jwt.expire")
	if t != "" {
		val, err := time.ParseDuration(t)
		if err != nil {
			return defaultExpireTime
		}
		return val
	}
	return defaultExpireTime
}
