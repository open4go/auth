package auth

import (
	"context"
	"strconv"
)

// CanAccess 是否允许访问
func CanAccess(ctx context.Context, keyForRoles string, path string, pathAccess string) bool {
	roles, err := RDB.SMembers(ctx, keyForRoles).Result()
	if err != nil {
		return false
	}

	// 仅进行路径的请求访问权限校验
	//pathAccess := AccessKeyPrefix + "_" + accountID + "_" + "path_access"
	for _, role := range roles {
		pathWithRole := path + "_" + role
		val, err := RDB.HGet(ctx, pathAccess, pathWithRole).Result()
		if err != nil {
			// 可以忽略该日志
			// 一般情况下仅角色匹配到path即可访问
			// 其他角色大部分会走该逻辑
			continue
		}
		// is true
		// 如果有一个角色是true 则代表其可以访问
		boolValue, err := strconv.ParseBool(val)
		if err != nil {
			// 可以忽略该日志
			// 一般情况下仅角色匹配到path即可访问
			continue
		}

		if boolValue {
			return true
		}
	}
	return false
}

// CRUD 是否允许操作
func CRUD(ctx context.Context, path string, keyOperation string, method string) bool {
	pathWithMethod := path + "/" + method
	val, err := RDB.HGet(ctx, keyOperation, pathWithMethod).Result()
	if err != nil {
		// 可以忽略该日志
		// 一般情况下仅角色匹配到path即可访问
		// 其他角色大部分会走该逻辑
		return false
	}
	// is true
	// 如果有一个角色是true 则代表其可以访问
	boolValue, err := strconv.ParseBool(val)
	if err != nil {
		// 可以忽略该日志
		// 一般情况下仅角色匹配到path即可访问
		// 其他角色大部分会走该逻辑
		return false
	}
	return boolValue
}
