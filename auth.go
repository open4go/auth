package auth

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/r2day/collections"
	log "github.com/sirupsen/logrus"
)

// SimpleAuth 基本类型
type SimpleAuth struct {
	MaxAccessLevel uint     `json:"max_access_level"`
	MyRoles        []string `json:"my_roles"`
	DisplayToolBar int      `json:"display_tool_bar"`
	// 键管理
	Key BasicKey `json:"key"`
	// 应用列表
	Apps []*AppModel `json:"apps"`
	// 角色配置
	RoleParam RoleParams `json:"role_param"`

	Path2Roles map[string][]string `json:"path_2_roles"`

	// 获取接口列表（所有角色下的接口列表)
	ApiList map[string][]collections.APIInfo `json:"api_list"`
	// 每个接口增删改查的权限
	Op []PermissionsModel `json:"op"`
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
	// 类型 hash 是否隐藏
	Hide string `json:"hide"`
	// path_access
	PathAccess string `json:"path_access"`
	//	role2paths
	Role2Paths string `json:"role2paths"`
	//	role2set4paths 集合类型
	// key 存储在 role2paths 的值中
	Role2Set4Paths string `json:"role2set4paths"`
}

// RoleParams 角色参数
type RoleParams struct {
	// 顶部工具栏,创建、导出、上传
	ToolBar int `json:"tool_bar"`
	// 最大可入级别 1~9
	MaxAccessLevel uint `json:"max_access_level"`
	// 角色名称列表
	// 用户登陆后，需要显示当前自己的角色
	RoleNameList []string `json:"role_name_list"`
	// 权限列表
	Permissions []PermissionsModel `json:"permissions"`
}

const (
	authPrefixKey = "auth_basic_"
)

// NewRBAM 新的角色验证模型
func NewRBAM() *SimpleAuth {
	return &SimpleAuth{
		MaxAccessLevel: 0,
		MyRoles:        make([]string, 0),
		DisplayToolBar: 0,
		Apps:           make([]*AppModel, 0),
		Path2Roles:     make(map[string][]string, 0),
		ApiList:        make(map[string][]collections.APIInfo, 0),
		Op:             make([]PermissionsModel, 0),
	}
}

// LoadConfig 加载配置
func (a *SimpleAuth) LoadConfig() {

}

func (a *SimpleAuth) HideMe(ctx context.Context, path string) bool {
	needToHideMe, err := RDB.HGet(ctx, a.Key.Hide, path).Result()
	if err != nil {
		return false
	}

	if needToHideMe == "1" {
		return true
	}
	return false
}

func (a *SimpleAuth) GetAllowPaths(ctx context.Context) []string {
	paths := make([]string, 0)
	myRoles := a.GetMyRoles(ctx)

	for _, role := range myRoles {
		tmpPaths := make([]string, 0)
		// user_1 是hash key，username 是字段名, tizi365是字段值
		secondKey, err := RDB.HGet(ctx, a.Key.Role2Paths, role).Result()
		if err != nil {
			log.WithField("secondKey", secondKey).Error(err)
			continue
		}
		tmpPaths, err = RDB.SMembers(ctx, secondKey).Result()
		if err != nil {
			log.WithField("tmpPaths secondKey", secondKey).Error(err)
			continue
		}
		paths = append(paths, tmpPaths...)
	}
	return paths
}

func (a *SimpleAuth) GetMyRoles(ctx context.Context) []string {
	roles, err := RDB.HGetAll(ctx, a.Key.Role2Paths).Result()
	if err != nil {
		log.Error(err)
		return a.MyRoles
	}
	for role, _ := range roles {
		a.MyRoles = append(a.MyRoles, role)
	}
	return a.MyRoles
}

// BindKey 绑定用户相关key
func (a *SimpleAuth) BindKey(accountID string) *SimpleAuth {
	keyPrefix := authPrefixKey + "_" + accountID
	a.Key = BasicKey{
		Keys:           keyPrefix + "_" + "keys",
		Roles:          keyPrefix + "_" + "roles",
		Operation:      keyPrefix + "_" + "operations",
		Path2Name:      keyPrefix + "_" + "path2name",
		Hide:           keyPrefix + "_" + "hide",
		PathAccess:     keyPrefix + "_" + "path_access",
		Role2Paths:     keyPrefix + "_" + "role2paths",
		Role2Set4Paths: keyPrefix + "_" + "role2set4paths",
	}

	// TODO 每一次操作都会更新expire，即当用户有操作行为则会延长过期时间
	err := a.ExpireSet(context.TODO())
	if err != nil {
		log.Error(err)
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
	err = RDB.SAdd(ctx, a.Key.Keys, a.Key.Role2Set4Paths).Err()
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
func (a *SimpleAuth) SignIn(ctx context.Context) error {
	err := a.recordKeys(ctx)
	if err != nil {
		return err
	}

	err = a.allowAccess(ctx, a.Path2Roles)
	if err != nil {
		return err
	}
	// 设置所有key过期时间

	return nil
}

// LoadRoles 加载角色
// 客户自行实现角色与账号的关联
// 角色信息会被加载到redis中
func (a *SimpleAuth) LoadRoles(ctx context.Context, roles []*RoleModel,
	apiListMap map[string][]collections.APIInfo) *SimpleAuth {
	a.ApiList = apiListMap
	// 显示工具栏: 创建、导出、上传
	//toolBar := 0
	rolesName := make([]string, 0)
	// 权限列表
	permissions := make([]PermissionsModel, 0)

	for _, role := range roles {
		log.WithField("apiInfo.Name", role.Name).Info("-------roles----")
		// 角色状态不可用
		if !role.Meta.Status {
			break
		}

		a.MyRoles = append(a.MyRoles, role.Name)
		// 选择最大的权限，决定是否展示状态栏
		// 权限越大，展示的功能越多
		if role.Toolbar > a.DisplayToolBar {
			a.DisplayToolBar = role.Toolbar
		}

		// 将所有角色下的接口的权限管理进行统一管理
		permissions = append(permissions, role.Permissions...)

		// 加入key
		// 以便退出登陆后删除
		rolesName = append(rolesName, role.Name)

		// 遍历寻找最大的用户角色等级
		// 以用户的最高权限进行登陆
		if role.Meta.AccessLevel > a.MaxAccessLevel {
			a.MaxAccessLevel = role.Meta.AccessLevel
		}

		// 将处理好的角色名称也加入到缓存中
		// 使用角色id 避免用户输入特殊字符无法作为redis key
		err := RDB.SAdd(ctx, a.Key.Roles, role.ID.Hex()).Err()
		if err != nil {
			log.WithField("roleName", role.Name).Error(err)
			continue
		}

		err = a.SetAccess(ctx, a.ApiList[role.ID.Hex()], role.ID.Hex())
		if err != nil {
			log.WithField("roleName", role.Name).Error(err)
			continue
		}
		log.WithField("apiInfo.Name", role.Name).Info("-------roles--done--")
	}

	// 设置permission
	err := a.setPermissions(ctx, permissions)
	if err != nil {
		return a
	}
	rp := RoleParams{
		ToolBar:        a.DisplayToolBar,
		MaxAccessLevel: a.MaxAccessLevel,
		RoleNameList:   rolesName,
		Permissions:    permissions,
	}
	a.RoleParam = rp

	return a
}

// Verify 加载应用
// 客户自行实现角色与账号的关联
// 角色信息会被加载到redis中
func (a *SimpleAuth) Verify(ctx context.Context, path string, method string) int {
	roles := a.GetMyRoles(ctx)
	// 检测角色是否有权限
	isAccess := CanAccess(ctx, roles, path, a.Key.PathAccess)
	if !isAccess {
		return http.StatusForbidden
	}

	// 检测账号是否有操作权限
	isCanDo := CanDo(ctx, path, a.Key.Operation, method)
	if !isCanDo {
		return http.StatusMethodNotAllowed
	}
	return http.StatusOK
}

// LoadApps 加载应用
// 客户自行实现角色与账号的关联
// 角色信息会被加载到redis中
//func (a *SimpleAuth) LoadApps(ctx context.Context, apps []*AppModel) error {
//	a.Apps = apps
//	return nil
//}

func (a *SimpleAuth) setPermissions(ctx context.Context, permissions []PermissionsModel) error {
	err := operatingAuthority(ctx, a.Key.Operation, permissions)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// SetAccess 返回目录列表
// 管理台根据返回的数据决定是否显示在导航栏
func (a *SimpleAuth) SetAccess(ctx context.Context, apiList []collections.APIInfo, roleID string) error {
	for _, apiInfo := range apiList {
		// 默认是false
		// 如果是true则忽略本条规则
		if apiInfo.Disable {
			continue
		}

		err := RDB.HSet(ctx, a.Key.Path2Name, apiInfo.Path, apiInfo.Name).Err()
		if err != nil {
			log.WithField("apiInfo.Name", apiInfo.Name).Error(err)
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

		a.Path2Roles[apiInfo.Path] = append(a.Path2Roles[apiInfo.Path], roleID)
		// 如果开启
		if apiInfo.CanViewDetail {
			pathForDetail := apiInfo.Path + "/:_id"
			a.Path2Roles[pathForDetail] = append(a.Path2Roles[pathForDetail], roleID)
		}
	}
	return nil
}

// AllowAccess 设定用户对api的最终可访问信息
// 仅当其设定了key value后才能进行访问
// 中间件会检测redis中是否设定了该key
// 与角色绑定
func (a *SimpleAuth) allowAccess(ctx context.Context, path2roles map[string][]string) error {
	// 定义角色-> 路径列表
	roles2paths := make(map[string][]string)

	// 加载默认可以开放的接口配置
	// config["admin"] = append(config["admin"], "/v1/auth/merchant/signin")
	// user_1 是hash key，username 是字段名, 是字段值
	// key := accessKeyPrefix + accountId

	for path, roles := range path2roles {
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

	//log.WithField("roles2paths", roles2paths).Debug("check the roles to paths")

	// 加载访问控制信息到redis中
	// 以便access及中间件完成check
	for role, paths := range roles2paths {
		// TODO 这里要记得删除key
		secondKey := a.Key.Role2Set4Paths + "_" + strconv.Itoa(int(time.Now().Unix()))
		// 将key 记录下来以便退出的时候进行删除
		// 将其本身也进行记录
		err := RDB.SAdd(ctx, a.Key.Keys, secondKey).Err()
		if err != nil {
			return err
		}
		// 将key 记录下来以便退出的时候进行删除
		// 将其本身也进行记录
		err = RDB.SAdd(ctx, secondKey, paths).Err()
		if err != nil {
			return err
		}
		err = RDB.HSet(ctx, a.Key.Role2Paths, role, secondKey).Err()
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

// ExpireSet 退出
func (a *SimpleAuth) ExpireSet(ctx context.Context) error {
	keys, err := RDB.SMembers(ctx, a.Key.Keys).Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err = RDB.Expire(ctx, key, getExpireTime()).Result()
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

// IsOnline 是否在线
func (a *SimpleAuth) IsOnline(ctx context.Context) (bool, error) {
	keys, err := RDB.SMembers(ctx, a.Key.Keys).Result()
	if err != nil {
		return false, err
	}
	return len(keys) > 0, nil
}
