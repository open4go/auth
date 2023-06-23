package auth

import (
	"context"

	"github.com/r2day/db"
	"github.com/r2day/middle"
	log "github.com/sirupsen/logrus"
)

func cleanAllKeysRelateToCurrentUserID(ctx context.Context, accountID string, keyNames string) error {

	// keyPrefix := middle.AccessKeyPrefix + "_" + accountID
	// keyNames := keyPrefix + "_" + "names"
	keys, err := RDB.SMembers(ctx, keyNames).Result()
	if err != nil {
		log.WithField("keyNames", keyNames).Error(err)
		return err
	}

	log.WithField("the_key", keyNames).
		WithField("keys", keys).Info("ready to delete keys")
	// 删除登陆后写入的缓存
	for _, key := range keys {
		log.WithField("fullKey", key).Info("ready to delete key")
		RDB.Del(ctx, key)
	}
	// RDB.Del(ctx, keyNames)

	// // 删除账号对应的信息
	// key := middle.AccessKeyPrefix + "_" + accountID
	// RDB.Del(ctx, key)


	return nil
}
