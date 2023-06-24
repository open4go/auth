package auth

import (
	"context"
	"encoding/json"
	"github.com/r2day/collections"
	"github.com/r2day/db"
	log "github.com/sirupsen/logrus"
)

// SimpleAuth 基本类型
type SimpleAuth struct {
	// 键管理
	Key BasicKey `json:"key"`
	// 应用列表
	Apps []*AppModel `json:"apps"`
	// 角色配置
	RoleParam RoleParams `json:"role_param"`

	Path2Roles map[string][]string `json:"path_2_roles"`
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
	// path_access
	PathAccess string `json:"path_access"`
	//	role2paths
	Role2Paths string `json:"role2paths"`
}

// RoleParams 角色参数
type RoleParams struct {
	// 顶部工具栏,创建、导出、上传
	ToolBar int `json:"tool_bar"`
	// 最大可入级别 1~9
	MaxAccessLevel int `json:"max_access_level"`
	// 角色名称列表
	// 用户登陆后，需要显示当前自己的角色
	RoleNameList []string `json:"role_name_list"`
	// 权限列表
	Permissions []PermissionsModel `json:"permissions"`
}

const (
	authPrefixKey = "auth_basic_"
)

// LoadConfig 加载配置
func (a *SimpleAuth) LoadConfig() {

}

// BindKey 绑定用户相关key
func (a *SimpleAuth) BindKey(accountID string) *SimpleAuth {
	keyPrefix := authPrefixKey + "_" + accountID
	a.Key = BasicKey{
		Keys:       keyPrefix + "_" + "keys",
		Roles:      keyPrefix + "_" + "roles",
		Operation:  keyPrefix + "_" + "operations",
		Path2Name:  keyPrefix + "_" + "path2name",
		Hide:       keyPrefix + "_" + "hide",
		PathAccess: keyPrefix + "_" + "path_access",
		Role2Paths: keyPrefix + "_" + "role2paths",
	}
	return a
}

// recordKeys 记录关联keys
func (a *SimpleAuth) recordKeys(ctx context.Context) error {
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
	err = RDB.SAdd(ctx, a.Key.Keys, a.Key.Roles).Err()
	if err != nil {
		return err
	}
	// 将key 记录下来以便退出的时候进行删除
	err = RDB.SAdd(ctx, a.Key.Keys, a.Key.Hide).Err()
	if err != nil {
		return err
	}

	// 将key 记录下来以便退出的时候进行删除
	err = RDB.SAdd(ctx, a.Key.Keys, a.Key.PathAccess).Err()
	if err != nil {
		return err
	}
	// 将key 记录下来以便退出的时候进行删除
	err = RDB.SAdd(ctx, a.Key.Keys, a.Key.Role2Paths).Err()
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
// 在密码与账号验证通过后调用该接口进行权限校验登陆
// 会将当前账号的权限与角色以及所有接口的具体操作权写入到redis中
func (a *SimpleAuth) SignIn(ctx context.Context, accountID string) error {
	err := a.recordKeys(ctx)
	if err != nil {
		return err
	}
	return nil
}

// LoadRoles 加载角色
// 客户自行实现角色与账号的关联
// 角色信息会被加载到redis中
func (a *SimpleAuth) LoadRoles(ctx context.Context, roles []*RoleModel) *SimpleAuth {
	// 显示工具栏: 创建、导出、上传
	toolBar := 0
	maxAccessLevel := 0
	rolesName := make([]string, 0)
	// 权限列表
	permissions := make([]PermissionsModel, 0)

	for _, role := range roles {

		// 角色状态不可用
		if !role.Meta.Status {
			break
		}

		// 选择最大的权限，决定是否展示状态栏
		// 权限越大，展示的功能越多
		if role.Toolbar > toolBar {
			toolBar = role.Toolbar
		}

		// 将所有角色下的接口的权限管理进行统一管理
		permissions = append(permissions, role.Permissions...)

		// 加入key
		// 以便退出登陆后删除
		rolesName = append(rolesName, role.Name)

		// 遍历寻找最大的用户角色等级
		if role.Meta.AccessLevel > uint(maxAccessLevel) {
			maxAccessLevel = int(role.Meta.AccessLevel)
		}

		// 将处理好的角色名称也加入到缓存中
		// 使用角色id 避免用户输入特殊字符无法作为redis key
		err := RDB.SAdd(ctx, a.Key.Roles, role.ID.Hex()).Err()
		if err != nil {
			continue
		}

		accessAPIList := make([]collections.APIInfo, 0)
		// 遍历所有授权的应用
		for _, app := range a.Apps {
			// 获得所有应用下的api列表
			accessAPIList = append(accessAPIList, app.AccessAPI...)
		}
		err = a.SetAccess(ctx, accessAPIList, role.ID.Hex())
		if err != nil {
			continue
		}
	}

	rp := RoleParams{
		ToolBar:        toolBar,
		MaxAccessLevel: maxAccessLevel,
		RoleNameList:   rolesName,
		Permissions:    permissions,
	}
	a.RoleParam = rp
	return a
}

// LoadApps 加载应用
// 客户自行实现角色与账号的关联
// 角色信息会被加载到redis中
func (a *SimpleAuth) LoadApps(ctx context.Context, apps []*AppModel) error {
	a.Apps = apps
	return nil
}

// SetAccess 返回目录列表
// 管理台根据返回的数据决定是否显示在导航栏
func (a *SimpleAuth) SetAccess(ctx context.Context, apiList []collections.APIInfo, roleID string) error {
	path2roles := make(map[string][]string, 0)
	for _, apiInfo := range apiList {
		// 默认是false
		// 如果是true则忽略本条规则
		if apiInfo.Disable {
			continue
		}

		err := RDB.HSet(ctx, a.Key.Path2Name, apiInfo.Path, apiInfo.Name).Err()
		if err != nil {
			continue
		}

		// 判断是否需要在导航menu中展示
		// 部分接口列表access和profile 是在个人中心展示的
		// 所以需要设置为true
		if apiInfo.HideOnSidebar {
			err = RDB.HSet(ctx, a.Key.Hide, apiInfo.Path, true).Err()
			if err != nil {
				continue
			}
		}

		path2roles[apiInfo.Path] = append(path2roles[apiInfo.Path], roleID)
		log.WithField("can_view_detail", apiInfo.CanViewDetail).Debug("check api info")
		// 如果开启
		if apiInfo.CanViewDetail {
			pathForDetail := apiInfo.Path + "/:_id"
			path2roles[pathForDetail] = append(path2roles[pathForDetail], roleID)
		}
	}
	a.Path2Roles = path2roles
	return nil
}

// Access 设定用户对api的最终可访问信息
// 仅当其设定了key value后才能进行访问
// 中间件会检测redis中是否设定了该key
// 与角色绑定
func (a *SimpleAuth) Access(ctx context.Context, accountID string, config map[string][]string) error {
	// 定义角色-> 路径列表
	roles2paths := make(map[string][]string)

	// 加载默认可以开放的接口配置
	// config["admin"] = append(config["admin"], "/v1/auth/merchant/signin")
	// user_1 是hash key，username 是字段名, 是字段值
	// key := accessKeyPrefix + accountId

	for path, roles := range config {
		for _, role := range roles {
			// api访问控制key
			pathWithRole := path + "_" + role
			err := RDB.HSet(ctx, a.Key.PathAccess, pathWithRole, true).Err()
			if err != nil {
				return err
			}
			roles2paths[role] = append(roles2paths[role], path)
		}
	}

	log.WithField("roles2paths", roles2paths).Debug("check the roles to paths")

	// 加载访问控制信息到redis中
	// 以便access及中间件完成check
	for role, paths := range roles2paths {
		pathsStr, err := json.Marshal(paths)
		if err != nil {
			log.Error(err)
			return err
		}
		err = db.RDB.HSet(ctx, a.Key.Role2Paths, role, pathsStr).Err()
		if err != nil {
			log.Error(err)
			return err
		}
	}

	// 判断是否存在如下key，否则报错
	// access_key_prefix_AC1657016915941396480_/v1/affiliate/membership_<role_id>_read true

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
