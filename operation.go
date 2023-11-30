package auth

import (
	"context"
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

		// 先重置删除旧的权限缓存
		err := resetOperatingAuthority(ctx, keyOperation, p.Path)
		if err != nil {
			return err
		}
		for _, requestParamWithMethod := range p.Operation {
			err = setOperatingAuthority(ctx, keyOperation, p.Path, requestParamWithMethod)
			if err != nil {
				return err
			}
		}
		//
		//// 写
		//pathWithWrite := p.Path + "/" + "POST"
		//err = setOperatingAuthority(ctx, keyOperation, pathWithWrite, p.Operation.Write)
		//if err != nil {
		//	return err
		//}
		//
		//// 改
		//pathWithUpdate := p.Path + "/:_id/" + "PUT"
		//err = setOperatingAuthority(ctx, keyOperation, pathWithUpdate, p.Operation.Update)
		//if err != nil {
		//	return err
		//}
		//// 详情
		//pathWithDetail := p.Path + "/:_id/" + "GET"
		//err = setOperatingAuthority(ctx, keyOperation, pathWithDetail, p.Operation.Detail)
		//if err != nil {
		//	return err
		//}
		//// 删除
		//pathWithDelete := p.Path + "/:_id/" + "DELETE"
		//err = setOperatingAuthority(ctx, keyOperation, pathWithDelete, p.Operation.Delete)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

// 根据用户操作的api path进行标记并写入数据库
func setOperatingAuthority(ctx context.Context, operatingAuthorityKey string, pathAndOperation string, requestParamWithMethod string) error {
	err := RDB.HSet(ctx, operatingAuthorityKey, pathAndOperation, requestParamWithMethod).Err()
	if err != nil {
		return err
	}
	return nil
}

// 根据用户操作的api path进行标记并写入数据库
func resetOperatingAuthority(ctx context.Context, operatingAuthorityKey string, pathAndOperation string) error {
	err := RDB.HDel(ctx, operatingAuthorityKey, pathAndOperation).Err()
	if err != nil {
		return err
	}
	return nil
}
