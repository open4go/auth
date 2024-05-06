package auth

import (
	"context"
	"strings"
)

var (
	map2op = map[string]string{
		"list":   "/GET",
		"write":  "/POST",
		"delete": "/:_id/DELETE",
		"update": "/:_id/PUT",
		"read":   "/:_id/GET",
	}
)

// operatingAuthority 操作权限 设定用户对api的最终可访问信息
// 仅当其设定了key value后才能进行访问
// 中间件会检测redis中是否设定了该key
// 与角色绑定
func operatingAuthority(ctx context.Context, keyOperation string, permissions []PermissionsModel) (err error) {
	// 加载默认可以开放的接口配置
	// config["admin"] = append(config["admin"], "/v1/auth/merchant/signin")
	// user_1 是hash key，username 是字段名, 是字段值
	// key := accessKeyPrefix + accountId

	for _, p := range permissions {
		allOp := strings.Join(p.Operation, ",")
		for op, v := range map2op {
			isTrue := false
			pathWithRead := p.Path + v
			if strings.Contains(allOp, op) {
				isTrue = true
			} else {
				isTrue = false
			}
			err = setOperatingAuthority(ctx, keyOperation, pathWithRead, isTrue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 根据用户操作的api path进行标记并写入数据库
func setOperatingAuthority(ctx context.Context, operatingAuthorityKey string, pathAndOperation string, val bool) error {
	err := GetRedisAuthHandler().HSet(ctx, operatingAuthorityKey, pathAndOperation, val).Err()
	if err != nil {
		return err
	}
	return nil
}
