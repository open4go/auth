package auth

import (
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

var RDB *redis.Client

func InitRedisDB(redisAddr string, db int, poolSize int) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:       redisAddr,
		DB:         db,
		PoolSize:   poolSize,
		MaxRetries: 3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := GetRedisAuthHandler().Ping(ctx).Err(); err != nil {
		log.WithFields(log.Fields{
			"redisAddr": redisAddr,
			"db":        db,
			"poolSize":  poolSize,
		}).WithError(err).Error("failed to connect to Redis")
		return err
	}

	log.WithFields(log.Fields{
		"redisAddr": redisAddr,
		"db":        db,
		"poolSize":  poolSize,
	}).Info("connected to Redis")

	return nil
}
