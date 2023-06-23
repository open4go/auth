package auth

import (
	"context"

	"github.com/r2day/collections"
	"github.com/r2day/collections/capp"
	"github.com/r2day/collections/crole"
	"github.com/r2day/db"
)

// 遍历账号中所拥有的所有Roles
// 1. 通过角色找到其应用列表，并且决定是否显示其在导航menu上或者是否禁用api
// 2. 通过permissions 找到每一个接口的具体操作方法，是否具备相关操作权限
func iterRoles(ctx context.Context, roles []*crole.Model,
	rolesName []string, maxAccessLevel int, myRolesKey string, keyPrefix string,
	path2roles map[string][]string) ([]string, []PermissionsModel, int, int) {

	// 显示工具栏: 创建、导出、上传
	toolBar := 0

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

		// 获取当前用户角色的应用列表
		appM := &capp.Model{}
		apps, err := appM.GetMany(ctx, role.Apps)
		if err != nil {
			continue
		}

		// 将处理好的角色名称也加入到缓存中
		// 使用角色id 避免用户输入特殊字符无法作为redis key
		err = db.RDB.SAdd(ctx, myRolesKey, role.ID.Hex()).Err()
		if err != nil {
			continue
		}

		accessAPIList := make([]collections.APIInfo, 0)

		// 遍历所有授权的应用
		for _, app := range apps {
			// 获得所有应用下的api列表
			accessAPIList = append(accessAPIList, app.AccessAPI...)
		}

		// setMenu(c, accessAPIList, keyPrefix, err, keyNames, path2roles, role)

	}

	return rolesName, permissions, maxAccessLevel, toolBar
}
