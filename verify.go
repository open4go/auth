package auth

import (
	"context"
	"github.com/open4go/log"
	"strconv"
)

func CanAccess(ctx context.Context, roles []string, path string, pathAccess string) bool {
	for _, role := range roles {
		pathWithRole := path + "_" + role
		val, err := GetRedisAuthHandler(ctx).HGet(ctx, pathAccess, pathWithRole).Result()
		if err != nil {
			logIgnorableWarning(ctx, "CanAccess", role, pathAccess, pathWithRole, err)
			continue
		}
		if boolValue, err := strconv.ParseBool(val); err == nil && boolValue {
			return true
		}
	}
	return false
}

func CanDo(ctx context.Context, path string, keyOperation string, method string) bool {
	pathWithMethod := path + "/" + method
	if val, err := GetRedisAuthHandler(ctx).HGet(ctx, keyOperation, pathWithMethod).Result(); err == nil {
		if boolValue, err := strconv.ParseBool(val); err == nil {
			return boolValue
		}
	}
	return false
}

func logIgnorableWarning(ctx context.Context, funcName string, role, pathAccess, pathWithRole string, err error) {
	log.Log(ctx).WithField("role", role).
		WithField("pathAccess", pathAccess).
		WithField("pathWithRole", pathWithRole).
		Warning(funcName, " - Ignorable: ", err)
}
