package auth

import (
	"context"
)

// SimpleAuth 基本类型
type SimpleAuth struct {
	// 键管理
	Key BasicKey `json:"key"`
}

// BasicKey 缓存键
type BasicKey struct {
	// 类型 set 存储当前登陆用户的 所有键
	// 以便当用户退出后进行统一删除
	Keys string `json:"keys"`
	// 类型 set 角色存储, 保存当前账号拥有的所有角色名称
	Roles string `json:"roles"`
	// 操作
	Operation string `json:"operation"`
	// 操作
	Path2Name string `json:"path_2_name"`
	// 是否隐藏
	Hide string `json:"hide"`
}

const (
	authPrefixKey = "auth_basic_"
	tokenKey      = "TOKEN_KEY"
)

// LoadConfig 加载配置
func (a *SimpleAuth) LoadConfig() {

}

// BindKey 绑定用户相关key
func (a *SimpleAuth) BindKey(accountID string) *SimpleAuth {
	keyPrefix := authPrefixKey + "_" + accountID
	a.Key = BasicKey{
		Keys:      keyPrefix + "_" + "keys",
		Roles:     keyPrefix + "_" + "roles",
		Operation: keyPrefix + "_" + "operations",
		Path2Name: keyPrefix + "_" + "path2name",
		Hide:      keyPrefix + "_" + "hide",
	}
	return a
}

// recordKeys 记录关联keys
func (a *SimpleAuth) recordKeys(ctx context.Context, accountID string) error {
	// 将key 记录下来以便退出的时候进行删除
	err := RDB.SAdd(ctx, a.Key.Keys, a.Key.Operation).Err()
	if err != nil {
		return err
	}
	// 将key 记录下来以便退出的时候进行删除
	err = RDB.SAdd(ctx, a.Key.Keys, a.Key.Path2Name).Err()
	if err != nil {
		return err
	}
	// 将key 记录下来以便退出的时候进行删除
	// 将其本身也进行记录
	err = RDB.SAdd(ctx, a.Key.Keys, a.Key.Keys).Err()
	if err != nil {
		return err
	}
	return nil
}

// SignIn 登陆
func (a *SimpleAuth) SignIn(ctx context.Context, accountID string) error {
	err := a.recordKeys(ctx, accountID)
	if err != nil {
		return err
	}
	return nil
}

// SignOut 退出
func (a *SimpleAuth) SignOut(ctx context.Context) error {
	keys, err := RDB.SMembers(ctx, a.Key.Keys).Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		RDB.Del(ctx, key)
	}
	return nil
}
